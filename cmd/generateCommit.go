package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

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

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	// 1. Point the client to Groq's OpenAI-compatible endpoint
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.groq.com/openai/v1"
	client := openai.NewClientWithConfig(config)

	// 2. Send the chat completion request
	ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "qwen/qwen3.8-27b",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fullPrompt,
				},
			},
			Temperature:         0.6,
			TopP:                0.95,
			MaxCompletionTokens: 2048,
		},
	)
	if err != nil {
		s.Stop()
		fmt.Printf("API Error: %v\n", err)
		return
	}
	if len(resp.Choices) > 0 {
		generatedMessage = strings.TrimSpace(resp.Choices[0].Message.Content)
	}

	if generatedMessage == "" {
		s.Stop()
		fmt.Println("Failed to generate commit message")
		return
	}

	s.Stop()
	fmt.Println("\nGenerated Commit Message >", generatedMessage)

	gitCommitCmd := exec.Command("git", "commit", "-m", generatedMessage)
	output, err := gitCommitCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error committing changes:", err)
		fmt.Println(string(output))
	} else {
		fmt.Println("Commit successful! Run git push")
	}
}
