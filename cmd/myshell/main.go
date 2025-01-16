package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

var builtinCmds = []string{"echo", "exit", "type", "pwd", "cd"}

func handleCommand(args []string) {
	cmd := strings.Trim(args[0], " ")
	args = args[1:]

	switch cmd {
	case "exit":
		if len(args) == 1 {
			exitCode, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Print(err)
			}
			os.Exit(exitCode)
		}
	case "echo":
		fmt.Println(strings.Join(args, " "))
		return
	case "type":
		cmdName := args[0]
		found := false
		for _, cmd := range builtinCmds {
			if cmd == cmdName {
				found = true
				break
			}
		}
		if found {
			fmt.Printf("%s is a shell builtin\n", cmdName)
		} else {
			filePath, found := findInPath(cmdName)

			if found {
				fmt.Printf("%s is %s\n", cmdName, filePath)
			} else {
				fmt.Printf("%s: not found\n", cmdName)
			}
		}
	case "pwd":
		pwd, _ := os.Getwd()
		fmt.Println(pwd)
	case "cd":
		dir := args[0]
		if dir == "~" {
			dir = os.Getenv("HOME")
		}
		err := os.Chdir(dir)
		if err != nil {
			fmt.Printf("cd: %s: No such file or directory\n", dir)
		}
	default:
		cmdToExecute := exec.Command(cmd, args...)
		cmdToExecute.Stderr = os.Stderr
		cmdToExecute.Stdout = os.Stdout
		err := cmdToExecute.Run()
		if err != nil {
			fmt.Printf("%s: command not found\n", cmd)
		}
	}
}

func findInPath(cmdName string) (string, bool) {
	pathDirs := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	for _, dir := range pathDirs {
		filePath := filepath.Join(dir, cmdName)
		if fileInfo, err := os.Stat(filePath); err == nil && !fileInfo.IsDir() && isExecutable(fileInfo.Mode()) {
			return filePath, true
		}
	}
	return "", false
}

func isExecutable(fileMode fs.FileMode) bool {
	// found how to do the check on stackoverflow
	return fileMode&0111 != 0
}

func handleSingleQuote(cmd string) []string {
	// Remove any trailing newlines or carriage returns
	ss := strings.Trim(cmd, "\r\n")
	var tokens []string
	var currentToken strings.Builder
	var inSingleQuote bool // tracks if we're inside single quotes
	var lastWasQuote bool  // tracks if we just finished a quoted section (for merging adjacent quotes)

	// First, extract the command (first word) using Fields to handle multiple spaces
	fields := strings.Fields(ss)
	if len(fields) == 0 {
		return nil
	}
	tokens = append(tokens, fields[0])

	// Remove the command portion from the string and any leading spaces
	ss = ss[len(fields[0]):]
	ss = strings.TrimSpace(ss)

	// Process the arguments character by character
	for i := 0; i < len(ss); i++ {
		ch := ss[i]

		if ch == '\'' {
			if !inSingleQuote {
				// Found opening quote
				inSingleQuote = true
				// Only append current token if:
				// 1. We have accumulated characters AND
				// 2. We're not continuing from a previous quote (for merging adjacent quotes)
				if currentToken.Len() > 0 && !lastWasQuote {
					tokens = append(tokens, currentToken.String())
					currentToken.Reset()
				}
			} else {
				// Found closing quote
				// Don't append the token yet - wait to see if another quote follows
				inSingleQuote = false
				lastWasQuote = true
			}
		} else {
			if inSingleQuote {
				// Inside quotes: preserve all characters literally, including spaces
				currentToken.WriteByte(ch)
			} else {
				// Outside quotes: handle spaces as token separators
				lastWasQuote = false // No longer in a position to merge quotes
				if ch != ' ' {
					// Collect non-space characters into current token
					currentToken.WriteByte(ch)
				} else if currentToken.Len() > 0 {
					// Space found: if we have a token, append it
					tokens = append(tokens, currentToken.String())
					currentToken.Reset()
				}
			}
		}
	}

	// Append any remaining characters as the final token
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		var cmdStr, err = bufio.NewReader(os.Stdin).ReadString('\n')
		// get rid of \n
		command := cmdStr[:len(cmdStr)-1]
		args := handleSingleQuote(command)

		if err != nil {
			fmt.Println("error ", err)
			os.Exit(1)
		}

		handleCommand(args)
	}
}
