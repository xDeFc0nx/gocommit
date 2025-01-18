package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var generatedMessage string

func generateCommit(cmd *cobra.Command, args []string) {
	url := "https://api-inference.huggingface.co/models/Qwen/Qwen2.5-Coder-32B-Instruct/v1/chat/completions"
	apiKey, err := LoadAPIKey()
	if err != nil {
		fmt.Printf(
			"apiKey is null did you Run gocommit set-api --key hf_yourapikeyhere\n",
		)
	}

	fullPrompt, err := prompt()
	if err != nil {
		fmt.Println(err)
		return
	}

	payload := map[string]interface{}{
		"model": "Qwen/Qwen2.5-Coder-32B-Instruct",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": fullPrompt,
			},
		},
		"max_tokens": 500,
		"stream":     false,
	}
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var apiResp APIResponse
	err = json.Unmarshal([]byte(body), &apiResp)
	if err != nil {
		fmt.Println("error parsing json response:", err)
		return
	}
	s.Stop()
	if len(apiResp.Choices) > 0 {
		generatedMessage = apiResp.Choices[0].Message.Content
		fmt.Println("Generated Commit Message:", generatedMessage)

	} else {
		fmt.Println("No message content found in the response.")
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Do you want to regenerate the commit message? (y/N): ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(strings.ToLower(userInput))
		if userInput == "y" || userInput == "yes" {
			fmt.Println("Regenerating commit message...")
			generateCommit(cmd, args)
		} else {
			if generatedMessage == "" {
				fmt.Println("Error: Generated commit message is empty.")
				return
			}
			gitCommitCmd := exec.Command("git", "commit", "-m", generatedMessage)
			output, err := gitCommitCmd.CombinedOutput()
			if err != nil {
				fmt.Println("Error committing changes:", err)
				fmt.Println(string(output))
			} else {
				fmt.Println("Commit successful! run git push")
			}
			break
		}
	}
}
