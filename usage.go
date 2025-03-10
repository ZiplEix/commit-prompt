package main

import "fmt"

func printUsage() {
	fmt.Println("Usage: commit-prompt")
	fmt.Println("Options:")
	fmt.Println("\t-h, --help    	Show this help message")
	fmt.Println("\t-v, --version 	Show the version")
	fmt.Println("\t--no-logs     	Do not include git logs in the prompt")
	fmt.Println("\t--max-logs    	Maximum number of logs to include in the prompt (0 for all logs)")
	fmt.Println("\t--output-format	Output format (clipboard or stdout)")
	fmt.Println("\t--commit-format 	Commit message format (conventional or gitmoji)")
	fmt.Println("\t--custom-format 	Custom format string to specify the output format, will be written at the end of the prompt")
}
