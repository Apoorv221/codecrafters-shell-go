package main
import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)
var Cmds = make(map[string]func(args []string) error)
func handleExit(args []string) error {
	var (
		exitCode int
		err      error
	)
	if len(args) == 1 {
		exitCode, err = strconv.Atoi(args[0])
		if err != nil {
			return err
		}
	}
	os.Exit(exitCode)
	return nil
}
func handleEcho(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(os.Stdout)
		return nil
	}
	for i := 0; i < len(args)-1; i++ {
		fmt.Fprintf(os.Stdout, "%s ", args[i])
	}
	fmt.Fprintln(os.Stdout, args[len(args)-1])
	return nil
}
func locateCmd(cmd string) (string, bool) {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			parts := strings.Split(e.Name(), ".")
			name := parts[0]
			if cmd == name {
				return fmt.Sprintf("%s/%s", dir, name), true
			}
		}
	}
	return "", false
}
func handleType(args []string) error {
	if len(args) != 1 {
		return nil
	}
	cmd := args[0]
	if _, ok := Cmds[cmd]; ok {
		fmt.Fprintf(os.Stderr, "%s is a shell builtin\n", cmd)
		return nil
	}
	if path, ok := locateCmd(cmd); ok {
		fmt.Fprintf(os.Stdout, "%s is %s\n", cmd, path)
		return nil
	}
	fmt.Fprintf(os.Stderr, "%s: not found\n", cmd)
	return nil
}
func handlePwd(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, dir)
	return nil
}
func handleCd(args []string) error {
	if len(args) == 0 {
		return nil
	}
	dir := args[0]
	if dir == "~" {
		dir = os.Getenv("HOME")
	}
	err := os.Chdir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", dir)
	}
	return nil
}
func main() {
	Cmds["exit"] = handleExit
	Cmds["echo"] = handleEcho
	Cmds["type"] = handleType
	Cmds["pwd"] = handlePwd
	Cmds["cd"] = handleCd
	for {
		fmt.Fprint(os.Stdout, "$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		s := strings.Trim(input, "\r\n")
		var tokens []string
		for {
			start := strings.Index(s, "'")
			if start == -1 {
				tokens = append(tokens, strings.Fields(s)...)
				break
			}
			tokens = append(tokens, strings.Fields(s[:start])...)
			s = s[start+1:]
			end := strings.Index(s, "'")
			token := s[:end]
			tokens = append(tokens, token)
			s = s[end+1:]
		}
		cmd := strings.ToLower(tokens[0])
		var args []string
		if len(tokens) > 1 {
			args = tokens[1:]
		}
		if fn, ok := Cmds[cmd]; ok {
			err := fn(args)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		} else if path, ok := locateCmd(cmd); ok {
			c := exec.Command(path, args...)
			o, err := c.Output()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Fprint(os.Stdout, string(o))
			}
		} else {
			fmt.Fprintf(os.Stdout, "%s: command not found\n", cmd)
		}
	}
}