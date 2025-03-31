package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

var maxLength = 72

func prompt() (string, error) {
	commitTypes := `
Choose a type from the list below that best describes the git diff:
- docs: Documentation only changes
- style: Changes that do not affect the meaning of the code (white-space, formatting, etc.)
- refactor: A code change that neither fixes a bug nor adds a feature
- perf: A code change that improves performance
- test: Adding missing tests or correcting existing tests
- build: Changes that affect the build system or external dependencies
- ci: Changes to CI configuration files and scripts
- chore: Changes that don't modify src or test files
- revert: Reverts a previous commit
- feat: A new feature
- fix: A bug fix
`
	generatePrompt := fmt.Sprintf(`
Generate a concise git commit message written in present tense for the following code diff with the specifications below:
- Commit message must be a maximum of ${maxLength} characters.
- Exclude unnecessary details like translations.
- don't wrap it in bash or anything.
- don't include reasoning 
- only include the commit message
- make sure it has the commit types.
- Your response will be passed directly into git commit.
`, maxLength)

	gh := exec.Command("git", "log", "--oneline")
	var ghb bytes.Buffer
	gh.Stdout = &ghb
	err := gh.Run()
	if err != nil {
		return "", fmt.Errorf("error running git log: %w", err)
	}

	gitHistory := content{files: gh.String()}

	gitcommand := exec.Command("git", "diff", "--cached")
	var out bytes.Buffer
	gitcommand.Stdout = &out
	err = gitcommand.Run()
	if err != nil {
		return "", fmt.Errorf("error running git diff: %w", err)
	}

	gitDiff := content{files: out.String()}

	if gitDiff.files == "" {
		return "No changes detected in git diff.", nil
	}

	fullPrompt := fmt.Sprintf(
		commitTypes,
		"%s\n%s\n\nGit History:\n%s",
		gitHistory.files,
		"%s\n%s\n\nCode Diff:\n%s",
		gitDiff.files,
		generatePrompt,

	)

	return fullPrompt, nil
}
