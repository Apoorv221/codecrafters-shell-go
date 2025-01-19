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
		// input = strings.TrimRight(input, "\n")
		// tokenizedInput := strings.Split(input, " ")
		
		// Tokenizing the input.
		s := strings.Trim(input, "\r\n")
		var tokenizedInput []string
		for {
			start := strings.Index(s, "'")
			if start == -1 {
				tokenizedInput = append(tokenizedInput, strings.Split(s, " ")...)
				break
			}
			tokenizedInput = append(tokenizedInput, strings.Fields(s[:start])...)
			s = s[start+1:]
			end := strings.Index(s, "'")
			token := s[:end]
			tokenizedInput = append(tokenizedInput, token)
			s = s[end+1:]
		}
		cmd := strings.ToLower(tokenizedInput[0]) 
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
		// Join parameters with a single space between them
		output := strings.Join(params, " ")

		// Collapse multiple spaces into a single space
		output = strings.Join(strings.Fields(output), " ")
	
		// Print the result
		fmt.Fprintf(os.Stdout, "%s\n", output)
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