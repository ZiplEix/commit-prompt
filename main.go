package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
)

func getModifiedFiles() ([]string, error) {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return nil, fmt.Errorf("no git repository found")
	}

	cmd := exec.Command("git", "diff", "--cached", "--name-only")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run git diff command: %w", err)
	}

	files := strings.Split(strings.TrimSpace(out.String()), "\n")

	if len(files) < 1 || files[0] == "" {
		return nil, nil
	}

	return files, nil
}

func getGitDiff() (string, error) {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return "", fmt.Errorf("no git repository found")
	}

	cmd := exec.Command("git", "diff", "--cached")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run git diff command: %w", err)
	}

	return out.String(), nil
}

func buildPrompt(modifiedFiles []string, gitDiff string) (string, error) {
	if len(modifiedFiles) < 1 {
		return "", fmt.Errorf("no modified file detected")
	}

	prompt := "Write me a commit message following the conventional commit format for this pending changes:\n\n"

	prompt += "Here are my modified files:\n"
	for _, file := range modifiedFiles {
		prompt += fmt.Sprintf("- %s\n", file)
		fileContent, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("failed to read file content: %w", err)
		}
		prompt += fmt.Sprintf("%s\n\n", fileContent)
	}

	prompt += "and here is the result of the command git diff --cached:\n"
	prompt += gitDiff

	prompt += "\n\nYour answer should only contain the commit message and the the body of the commit without enything else."

	return prompt, nil
}

func main() {
	modifiedFile, err := getModifiedFiles()
	if err != nil {
		fmt.Println(err)
		return
	}

	gitDiff, err := getGitDiff()
	if err != nil {
		fmt.Println(err)
		return
	}

	prompt, err := buildPrompt(modifiedFile, gitDiff)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = clipboard.WriteAll(prompt)
	if err != nil {
		log.Fatalf("❌ Erreur lors de la copie dans le presse-papiers : %v", err)
	}

	fmt.Println("✅ Le prompt a été copié dans le presse-papiers. Vous pouvez maintenant le coller.")
}
