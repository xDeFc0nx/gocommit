package cmd

import (
	"fmt"
	"os"
	"strings"
)

func LoadAPIKey() (string, error) {
	homeDir, _ := os.UserHomeDir()
	filePath := homeDir + "/.gocommit"

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf(
			"could not read .gocommit file did you Run gocommit set-api --key yourapikeyhere:  %w",
			err,
		)
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "API_KEY=") {
			return strings.TrimPrefix(line, "API_KEY="), nil
		}
	}

	return "", fmt.Errorf(
		"API_KEY not found in .gocommit file, Run gocommit set-api --key yourapikeyhere",
	)
}
