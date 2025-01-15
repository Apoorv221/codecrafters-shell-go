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
		tokenizedInput := strings.Split(input, " ")
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
	err := os.Chdir(params[0])
	if err!=nil{
		fmt.Fprintln(os.Stdout,"cd: %s: No such file or directory",params)
	}
	
}