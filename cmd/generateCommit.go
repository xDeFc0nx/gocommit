package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

type content struct {
	files string
}
type APIResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

var generatedMessage string

func generateCommit(cmd *cobra.Command, args []string) {
	apiKey, err := LoadAPIKey()
	if err != nil {
		fmt.Println("API key not found. Configure using: gocommit set-api --key YOUR_API_KEY")
		return
	}

	fullPrompt, err := prompt()
	if err != nil {
		fmt.Println("Error getting diff:", err)
		return
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s",
		apiKey,
	)

	payload := map[string]any{
		"contents": []any{
			map[string]any{
				"parts": []any{
					map[string]any{
						"text": fullPrompt,
					},
				},
			},
		},
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("API Error (%d): %s\n", resp.StatusCode, string(body))
		return
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	if len(apiResp.Candidates) > 0 &&
		len(apiResp.Candidates[0].Content.Parts) > 0 {
		generatedMessage = apiResp.Candidates[0].Content.Parts[0].Text
		generatedMessage = strings.TrimSpace(generatedMessage)
	}

	if generatedMessage == "" {
		fmt.Println("Failed to generate commit message")
		return
	}

	fmt.Println("\nGenerated Commit Message >", generatedMessage)
	gitCommitCmd := exec.Command("git", "commit", "-m", generatedMessage)
	output, err := gitCommitCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error committing changes:", err)
		fmt.Println(string(output))
	} else {
		fmt.Println("Commit successful! run git push")
	}
}
