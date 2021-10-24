package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Op int

type Command struct {
	op    Op
	key   string
	value string
}

const (
	Get Op = iota
	Set
)

var store map[string]string

func HandleInt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nExiting. Bye!")
		os.Exit(0)
	}()
}

func ParseCommand(cmd string) (*Command, error) {
	cmd = strings.TrimSpace(cmd)
	i := strings.Index(cmd, " ")
	if i < 0 {
		return nil, fmt.Errorf("Invalid command: not enough arguments")
	}

	op := strings.ToLower(cmd[0:i])
	args := cmd[i+1:]
	command := &Command{}
	if op == "get" {
		command.op = Get
		command.key = args
	} else if op == "set" {
		command.op = Set
		parts := strings.Split(args, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid command: must provide 'key=value' for set")
		}
		command.key = parts[0]
		command.value = parts[1]
	} else {
		return nil, fmt.Errorf("Invalid command: unknown operation")
	}
	return command, nil
}

func ExecuteCommand(cmd *Command) (string, error) {
	switch op := cmd.op; op {
	case Get:
		if val, ok := store[cmd.key]; ok {
			return val, nil
		} else {
			return "", fmt.Errorf("Key '%s' not found", cmd.key)
		}
	case Set:
		store[cmd.key] = cmd.value
		return fmt.Sprintf("Set key '%s' to '%s'", cmd.key, cmd.value), nil
	default:
		return "", fmt.Errorf("Invalid Operation: %v", op)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	store = make(map[string]string)
	HandleInt()

	for {
		fmt.Printf("dkvs> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input")
			continue
		}
		cmd, err := ParseCommand(line)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		result, err := ExecuteCommand(cmd)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(result)
	}
}
