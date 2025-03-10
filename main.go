package main

import (
	"bytes"
	"flag"
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

	fmt.Println("🔄 Getting modified files...")

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

	fmt.Println("🔄 Getting git diff...")

	cmd := exec.Command("git", "diff", "--cached")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run git diff command: %w", err)
	}

	return out.String(), nil
}

func getGitLogs(maxLogs int) (string, error) {
	// fmt.Println("🔄 Récupération des logs git...")
	fmt.Println("🔄 Getting git logs...")

	var cmd *exec.Cmd
	if maxLogs <= 0 {
		cmd = exec.Command("git", "log")
	} else {
		cmd = exec.Command("git", "log", fmt.Sprintf("-%d", maxLogs))
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run git log command: %w", err)
	}

	return out.String(), nil
}

func buildPrompt(modifiedFiles []string, gitDiff string, log string, maxLogs int) (string, error) {
	if len(modifiedFiles) < 1 {
		return "", fmt.Errorf("no modified file detected")
	}

	fmt.Println("🔄 Building prompt...")

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

	if log != "" {
		if maxLogs <= 0 {
			prompt += "\nFor reference, here are the previous commits:\n"
		} else {
			prompt += fmt.Sprintf("\nFor reference, here are the previous %d commits:\n", maxLogs)
		}
		prompt += log
	}

	prompt += "\n\nYour answer should only contain the commit message and the the body of the commit without enything else."

	return prompt, nil
}

func generatePrompt(includeLogs bool, maxLogs int) {
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

	logs := ""
	if includeLogs {
		logs, err = getGitLogs(maxLogs)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	prompt, err := buildPrompt(modifiedFile, gitDiff, logs, maxLogs)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = clipboard.WriteAll(prompt)
	if err != nil {
		log.Fatalf("❌ Error while copying to clipboard: %v", err)
	}

	fmt.Println("✅ The prompt has been copied to the clipboard. You can now paste it.")
}

func main() {
	h := flag.Bool("h", false, "Show help message")
	help := flag.Bool("help", false, "Show help message")
	v := flag.Bool("v", false, "Show version")
	version := flag.Bool("version", false, "Show version")
	noLogs := flag.Bool("no-logs", false, "Do not include git logs in the prompt")
	maxLogs := flag.Int("max-logs", 0, "Maximum number of logs to include in the prompt (0 for all logs)")

	flag.Parse()

	if *h || *help {
		printUsage()
		return
	}

	if *v || *version {
		printVersion()
		return
	}

	generatePrompt(!(*noLogs), *maxLogs)
}
