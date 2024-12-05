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
		argStr := strings.Join(args, " ")
		if argStr[0] == '\'' && argStr[len(argStr)-1] == '\'' {
			fmt.Println(argStr[1 : len(argStr)-1])
			return
		}
		fmt.Println(strings.Join(args, " "))
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
