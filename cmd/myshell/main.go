package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func handleCommand(command string) {
	if command == "exit 0" {
		os.Exit(0)
	}

	split := strings.Split(command, " ")
	switch split[0] {
	case "echo":
		fmt.Println(strings.Join(split[1:], " "))
	default:
		fmt.Println(command + ": command not found")
	}
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
