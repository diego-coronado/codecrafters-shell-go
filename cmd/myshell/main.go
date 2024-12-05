package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	var command, err = bufio.NewReader(os.Stdin).ReadString('\n')

	if err != nil {
		fmt.Println("error ", err)
		os.Exit(1)
	}

	fmt.Println(command[:len(command)-1] + ": command not found")
}
