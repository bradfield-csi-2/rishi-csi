package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
var storeFile *os.File

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
			return "", fmt.Errorf("key '%s' not found", cmd.key)
		}
	case Set:
		store[cmd.key] = cmd.value
		jsonData, err := json.Marshal(store)
		if err != nil {
			return "", fmt.Errorf("Error marshaling JSON")
		}
		storeFile.WriteAt(jsonData, 0)
		return fmt.Sprintf("Set key '%s' to '%s'", cmd.key, cmd.value), nil
	default:
		return "", fmt.Errorf("Invalid Operation: %v", op)
	}
}

func SetupStore() *os.File {
	store = make(map[string]string)
	f, err := os.OpenFile("store.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening data store")
		os.Exit(1)
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading data store")
		os.Exit(1)
	}
	var out map[string]interface{}
	err = json.Unmarshal(bytes, &out)
	if err != nil {
		fmt.Println("Error unmarshaling JSON data from store")
		os.Exit(1)
	}

	for k, v := range out {
		store[k] = v.(string)
	}
	return f
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	HandleInt()
	storeFile = SetupStore()
	defer storeFile.Close()

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
