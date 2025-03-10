# Commit Prompt

Commit Prompt is a CLI tool written in Go that generates a detailed prompt to help write commit messages following the Conventional Commits format.

## Features

- Identifies modified files in the Git index.
- Extracts the diff of staged changes.
- Builds a prompt containing the modified files, their contents, and the `git diff --cached` output.
- Automatically copies the prompt to the clipboard for easy usage.

## Installation

Install the tool using `go install`:

```sh
go install github.com/ZiplEix/commit-prompt@latest
```

Ensure that `$GOPATH/bin` is in your system's PATH so you can run the command globally.

## Usage

In a Git repository, run:

```sh
commit-prompt
```

The generated prompt will be automatically copied to the clipboard. You can then paste it into an AI tool or use it to write a clear and convention-compliant commit message.

## Requirements

- Go 1.18+
- A Git repository with staged files (`git add` required)

## Example

If you have modified and staged two files, `main.go` and `utils.go`:

```sh
commit-prompt
```

This will copy the following message to the clipboard:

```
Write me a commit message following the conventional commit format for this pending changes:

Here are my modified files:
- main.go
<content of main.go>

- utils.go
<content of utils.go>

and here is the result of the command git diff --cached:
<git diff output>

Your answer should only contain the commit message and the body of the commit without anything else.
```

## Contributions

Contributions are welcome! Feel free to open an issue or a pull request on [GitHub](https://github.com/ZiplEix/commit-prompt).
