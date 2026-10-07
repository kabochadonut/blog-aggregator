package main

import (
	"blog-aggregator/internal/database"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		log.Fatalf("Usage: login <username>")
	}
	name := strings.ToLower(cmd.args[0])
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, name)
	if err != nil {
		log.Fatalf("User: '%s' doesn't exist\n", name)
	}
	s.config.SetUser(user.Name)
	fmt.Printf("User: '%s' has been set\n", user.Name)
	return nil
}

// Register new user
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		name := strings.ToLower(cmd.args[0])
		ctx := context.Background()
		_, err := s.db.GetUser(ctx, name)
		if err == nil {
			log.Fatalf("User exists")
		}
		params := database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
		}
		user, err := s.db.CreateUser(ctx, params)
		if err != nil {
			return fmt.Errorf("Error creating user: %v", err)
		}
		s.config.SetUser(user.Name)
		fmt.Printf("User '%s' was registered\n", name)
	} else {
		return errors.New("Usage: register <username>")
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		log.Fatalf("Failed to reset users: %v", err)
	}
	return err
}

func handlerGetUsers(s *state, cmd command) error {
	ctx := context.Background()
	usernames, err := s.db.GetUsers(ctx)
	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}
	for _, name := range usernames {
		if name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", name)
		} else {
			fmt.Printf("* %s\n", name)
		}
	}
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	if fn, ok := c.commands[cmd.name]; ok {
		err := fn(s, cmd)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("Command not found: %s", cmd.name)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}
