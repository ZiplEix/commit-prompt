package main

import "fmt"

func printUsage() {
	fmt.Println("Usage: commit-prompt")
	fmt.Println("Options:")
	fmt.Println("\t-h, --help    Show this help message")
	fmt.Println("\t-v, --version Show the version")
	fmt.Println("\t--no-logs     Do not include git logs in the prompt")
	fmt.Println("\t--max-logs    Maximum number of logs to include in the prompt (0 for all logs)")
}
