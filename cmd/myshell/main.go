package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

var builtinCmds = []string{"echo", "exit", "type"}

func handleCommand(command string) {
	args := strings.Split(command, " ")
	cmd := args[0]
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
	case "type":
		check := args[0]
		found := false
		for _, cmd := range builtinCmds {
			if cmd == check {
				found = true
				break
			}
		}
		if found {
			fmt.Printf("%s is a shell builtin\n", check)
		} else {
			pathDirs := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
			for _, dir := range pathDirs {
				filePath := filepath.Join(dir, check)
				if fileInfo, err := os.Stat(filePath); err == nil && !fileInfo.IsDir() && isExecutable(fileInfo.Mode()) {
					fmt.Printf("%s is %s\n", check, filePath)
					found = true
					break
				}
			}

			if !found {
				fmt.Printf("%s: not found\n", check)
			}
		}
	default:
		fmt.Println(command + ": command not found")
	}
}

func isExecutable(fileMode fs.FileMode) bool {
	// found how to do the check on stackoverflow
	return fileMode&0111 != 0
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		var cmdStr, err = bufio.NewReader(os.Stdin).ReadString('\n')
		// get rid of \n
		command := cmdStr[:len(cmdStr)-1]

		if err != nil {
			fmt.Println("error ", err)
			os.Exit(1)
		}

		handleCommand(command)
	}
}
