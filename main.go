package main

import (
	"blog-aggregator/internal/config"
	"fmt"
	"os"
)

func main() {
	cfg := config.Read()
	appState := state{
		config: &cfg,
	}
	commands := commands{
		commands: make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: <command> [args]")
		os.Exit(1)
	}
	cmdName := args[1]
	cmdArgs := args[2:]
	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}
	err := commands.run(&appState, cmd)
	if err != nil {
		fmt.Printf("Command: '%s' failed to run\n", cmdName)
		os.Exit(1)
	}
}
