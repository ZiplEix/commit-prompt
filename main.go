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

type parameters struct {
	includeLogs  bool
	maxLogs      int
	outputFormat string
	commitFormat string
	customFormat string
}

var params parameters

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

func buildPrompt(modifiedFiles []string, gitDiff string, log string) (string, error) {
	if len(modifiedFiles) < 1 {
		return "", fmt.Errorf("no modified file detected")
	}

	fmt.Println("🔄 Building prompt...")

	prompt := ""

	if params.customFormat == "" {
		prompt += fmt.Sprintf("Write me a commit message following the %s convention for this pending changes:\n\n", params.commitFormat)
	} else {
		prompt += fmt.Sprintf("Write me a commit message for this pending changes: The commit format should follow this rule: %s\n\n", params.customFormat)
	}

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
		if params.maxLogs <= 0 {
			prompt += "\nFor reference, here are the previous commits:\n"
		} else {
			prompt += fmt.Sprintf("\nFor reference, here are the previous %d commits:\n", params.maxLogs)
		}
		prompt += log
	}

	prompt += "\n\nYour answer should only contain the commit message and the the body of the commit without enything else."

	return prompt, nil
}

func generatePrompt() {
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
	if params.includeLogs {
		logs, err = getGitLogs(params.maxLogs)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	prompt, err := buildPrompt(modifiedFile, gitDiff, logs)
	if err != nil {
		fmt.Println(err)
		return
	}

	if params.outputFormat == "stdout" {
		fmt.Println(prompt)
		return
	} else if params.outputFormat == "clipboard" {
		err = clipboard.WriteAll(prompt)
		if err != nil {
			log.Fatalf("❌ Error while copying to clipboard: %v", err)
		}

		fmt.Println("✅ The prompt has been copied to the clipboard. You can now paste it.")
	}
}

func checkParams() error {
	if params.outputFormat != "clipboard" && params.outputFormat != "stdout" {
		fmt.Println("❌ Error: invalid output format")
		return fmt.Errorf("invalid output format")
	}

	if params.commitFormat != "conventional linux" && params.commitFormat != "gitmoji" {
		fmt.Printf("❌ Error: invalid commit format, have %s, but want %s or %s\n", params.commitFormat, "conventional", "gitmoji")
		return fmt.Errorf("invalid commit format")
	}

	return nil
}

func main() {
	h := flag.Bool("h", false, "Show help message")
	help := flag.Bool("help", false, "Show help message")
	v := flag.Bool("v", false, "Show version")
	version := flag.Bool("version", false, "Show version")
	noLogs := flag.Bool("no-logs", false, "Do not include git logs in the prompt")
	maxLogs := flag.Int("max-logs", 0, "Maximum number of logs to include in the prompt (0 for all logs)")
	outputFormat := flag.String("output-format", "clipboard", "Output format (clipboard, stdout)")
	commitFormat := flag.String("commit-format", "conventional linux", "Commit format (conventional linux, gitmoji)")
	customFormat := flag.String("custom-format", "", "Custom format string to specify the output format, wil be write at the end on the prompt")

	flag.Parse()

	if *h || *help {
		printUsage()
		return
	}

	if *v || *version {
		printVersion()
		return
	}

	params.includeLogs = !(*noLogs)
	params.maxLogs = *maxLogs
	params.outputFormat = *outputFormat
	params.commitFormat = *commitFormat
	params.customFormat = *customFormat

	err := checkParams()
	if err != nil {
		return
	}

	generatePrompt()
}
