package main
import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"os/exec"
)
var KnownCommands = map[string]int{"exit": 0, "echo": 1, "type": 2,"pwd" :3, "cd": 4}
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")
	for {
		// Uncomment this block to pass the first stage
		fmt.Fprint(os.Stdout, "$ ")
		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}
		// remove user enter
		input = strings.TrimRight(input, "\n")
		tokenizedInput := tokenizeInput(input)
		cmd := tokenizedInput[0]
		if fn, exists := KnownCommands[cmd]; !exists {
			
			resultCommand := exec.Command(cmd,tokenizedInput[1:]...)
			resultCommand.Stderr = os.Stderr
			resultCommand.Stdout = os.Stdout

			err := resultCommand.Run()
			if err != nil {
				fmt.Fprintf(os.Stdout,"%s: command not found\n",tokenizedInput[0])
			} 
		} else {
			switch fn {
			case 0:
				DoExit(tokenizedInput[1:])
			case 1:
				DoEcho(tokenizedInput[1:])
			case 2:
				DoType(tokenizedInput[1:])
			case 3: 
				DoPwd()
			case 4:
				DoCd(tokenizedInput[1:])
			}
		}
	}
}
func DoExit(params []string) {
	os.Exit(0)
}
func DoEcho(params []string) {
	output := strings.Join(params, " ")
	fmt.Fprintf(os.Stdout, "%v\n", output)
}
func DoType(params []string) {
	item := params[0]
	if _, exists := KnownCommands[item]; exists {
		class := "builtin"
		fmt.Fprintf(os.Stdout, "%v is a shell %v\n", item, class)
	} else {
		env := os.Getenv("PATH")
		paths := strings.Split(env, ":")
		for _, path := range paths {
			exec := path + "/" + item
			if _, err := os.Stat(exec); err == nil {
				fmt.Fprintf(os.Stdout, "%v is %v\n", item, exec)
				return
			}
		}
		fmt.Fprintf(os.Stdout, "%v not found\n", item)
	}
}

func DoPwd(){
	currentPath,err := os.Getwd(); 
	if err==nil{
		fmt.Fprintln(os.Stdout,currentPath)
		return
	}
}

func DoCd(params []string){
	if(params[0] == "~"){
		homePath,_ := os.UserHomeDir()
		os.Chdir(homePath)		
	}else{
		err := os.Chdir(params[0])
		if err!=nil{
			fmt.Fprintf(os.Stdout,"cd: %v: No such file or directory\n",params[0])
		}
	}
	
	
}
func tokenizeInput(input string) []string {
	var tokens []string
	var currentToken strings.Builder
	var inQuotes bool
	var quoteChar rune
	var escapeNext bool

	for _, char := range input {
		
		if escapeNext {
			currentToken.WriteRune(char)
			escapeNext = false
			continue
		}

		switch char {
		case '\'', '"':
			if inQuotes && char == quoteChar {
				// Close the quote
				inQuotes = false
			} else if !inQuotes {
				// Open a new quote
				inQuotes = true
				quoteChar = char
			} else {
				// Inside quotes, add the quote character
				currentToken.WriteRune(char)
			}
		case ' ':
			if inQuotes {
				// Inside quotes, spaces are part of the token
				currentToken.WriteRune(char)
			} else if currentToken.Len() > 0 {
				// Outside quotes, end of a token
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		case '\\':
			escapeNext = true
		default:
			// Normal characters
			currentToken.WriteRune(char)
		}
	}

	// Add the last token if any
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}