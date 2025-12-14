package main

import (
	"log"
	"os"

	"github.com/jacobovii/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	appState := &state{cfg: &cfg}

	commands := commands{
		registeredCommands: map[string]func(*state, command) error{},
	}
	commands.register("login", handleLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("not enough arguments provided")
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	err = commands.run(appState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
