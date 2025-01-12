package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func handleCommand (command string){
	fmt.Println(command + ": command not found")
	fmt.Println("$ exit 0")
	os.Exit(1)
}

func main() {
	// Uncomment this block to pass the first stage
	fmt.Fprint(os.Stdout, "$ ")

	for{
		reader := bufio.NewReader(os.Stdin)
		command, err := reader.ReadString('\n')

		if err != nil {
			os.Exit(1);
		}

		command = strings.TrimSuffix(command, "\n")
		handleCommand(command)
		fmt.Fprint(os.Stdout,"$ ")
	}	
}

