package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var apiKey string

var setAPIKeyCmd = &cobra.Command{
	Use:   "set-api",
	Short: "Set the API key for gocommit",
	Long:  `Set the API key used by gocommit to authenticate with external services.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := saveAPIKey(apiKey); err != nil {
			fmt.Println("Error saving API key:", err)
		} else {
			fmt.Println("API key saved successfully!")
		}
	},
}

func saveAPIKey(key string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".gocommit")

	file, err := os.OpenFile(
		configPath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0600,
	)
	if err != nil {
		return fmt.Errorf("failed to open .gocommit file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("API_KEY=%s\n", key))
	if err != nil {
		return fmt.Errorf("failed to write API key to file: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(setAPIKeyCmd)

	setAPIKeyCmd.Flags().
		StringVarP(&apiKey, "key", "k", "", "API key to save (required)")
	setAPIKeyCmd.MarkFlagRequired("key")
}
