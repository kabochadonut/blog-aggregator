package main

import _ "github.com/lib/pq"
import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"database/sql"
	"fmt"
	"log"
	"os"
)

func main() {
	cfg := config.Read()
	db, err := sql.Open("postgres", cfg.DBURL)
	dbQueries := database.New(db)
	appState := state{
		db:     dbQueries,
		config: &cfg,
	}

	commands := commands{
		commands: make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerGetUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", handlerAddFeed)
	commands.register("feeds", handlerFeeds)
	args := os.Args
	if len(args) == 1 {
		log.Fatalf("Usage: <command> [args]")
	}
	cmdName := args[1]
	cmdArgs := []string{}
	if len(args) > 2 {
		cmdArgs = args[2:]
	}
	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}
	err = commands.run(&appState, cmd)
	if err != nil {
		fmt.Printf("Command: '%s' failed to run\n", cmdName)
		log.Fatalf("Error: %v", err)
	}

}
