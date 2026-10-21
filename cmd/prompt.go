package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const maxTitleLength = 72

func prompt() (string, error) {
	diffCmd := exec.Command("git", "diff", "--cached")
	var diffOut bytes.Buffer
	diffCmd.Stdout = &diffOut
	if err := diffCmd.Run(); err != nil {
		return "", fmt.Errorf("error running git diff: %w", err)
	}

	diff := strings.TrimSpace(diffOut.String())
	if diff == "" {
		return "", errors.New("no staged changes found. Stage files with 'git add' before committing")
	}

	logCmd := exec.Command("git", "log", "-n", "5", "--oneline")
	var logOut bytes.Buffer
	logCmd.Stdout = &logOut
	_ = logCmd.Run() // non-fatal if this is a fresh repo with no commits yet
	recentCommits := strings.TrimSpace(logOut.String())

	promptTemplate := `You are an expert developer generating a Conventional Git Commit message based on a staged diff.

<instructions>
1. Format: <type>(<optional scope>): <imperative summary>
2. Allowed Types:
   - feat: A new feature
   - fix: A bug fix
   - refactor: Code restructuring without functional change
   - perf: Performance optimization
   - style: Formatting, missing semicolons, whitespace (no logic change)
   - test: Adding or updating tests
   - docs: Documentation changes only
   - build: Changes affecting build system or external dependencies
   - ci: Changes to CI workflows/scripts
   - chore: Tooling, configs, or maintenance tasks
3. Rules:
   - Must be written in the present imperative mood (e.g., "add feature", NOT "added feature" or "adds feature").
   - Total length of the first line MUST NOT exceed %d characters.
   - Do NOT capitalize the first letter after the colon.
   - Do NOT end the line with a period.
   - Infer the (<scope>) from modified file paths or modules if clear (e.g., "auth", "parser", "cli").
   - Output ONLY the raw commit message.
   - NO markdown fences (no ` + "```" + `), NO quotes, NO reasoning, NO explanation.
</instructions>

<recent_commits>
%s
</recent_commits>

<git_diff>
%s
</git_diff>
`

	if recentCommits == "" {
		recentCommits = "No prior commits (initial repository state)."
	}

	fullPrompt := fmt.Sprintf(promptTemplate, maxTitleLength, recentCommits, diff)
	return fullPrompt, nil
}
