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
}
func exitCommand(command string){
	os.Exit(0)
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
		
		if(command == "exit 0"){
			exitCommand(command)
		}
		handleCommand(command)
		
		fmt.Fprint(os.Stdout,"$ ")
	}	
}

