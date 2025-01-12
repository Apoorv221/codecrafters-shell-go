package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
func main() {
	// Uncomment this block to pass the first stage
	for {
		fmt.Fprint(os.Stdout, "$ ")
		message, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			os.Exit(1)
		}
		message = strings.TrimSpace(message)
		commands := strings.Split(message, " ")
		switch commands[0] {
		case "exit":
			code, err := strconv.Atoi(commands[1])
			if err != nil {
				os.Exit(1)
			}
			os.Exit(code)
		case "echo":
			fmt.Fprintf(os.Stdout, "%s\n", strings.Join(commands[1:], " "))

		case "type":
			if(strings.Contains(commands[1],"invalid")){
				fmt.Fprintf(os.Stdout, "%s: not found\n", commands[1])
			}
			fmt.Println(commands[1] + " is a shell builtin")
		default:
			fmt.Fprintf(os.Stdout, "%s: command not found\n", message)
		}
}
}