package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jacobovii/gator/internal/config"
	"github.com/jacobovii/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	appState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	commands := commands{
		registeredCommands: map[string]func(*state, command) error{},
	}
	commands.register("login", handleLogin)
	commands.register("register", handleRegister)
	commands.register("reset", handleReset)
	commands.register("users", handleGetUsers)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	err = commands.run(appState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
