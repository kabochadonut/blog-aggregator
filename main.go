package main

import _ "github.com/lib/pq"
import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"context"
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
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerListFeedFollows))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
