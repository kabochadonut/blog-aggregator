package main

import (
	"blog-aggregator/internal/database"
	"blog-aggregator/internal/rss"
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

func handlerAgg(s *state, cmd command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Usage: <command> [args]")
	}
	name := cmd.args[0]
	url := cmd.args[1]
	username := s.config.CurrentUserName
	user, err := s.db.GetUser(context.Background(), username)
	id := user.ID
	if err != nil {
		return fmt.Errorf("Failed to retrieve user: %v", err)
	}
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    id,
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Failed to create feed: %v", err)
	}
	fmt.Printf("Feed: '%s' at '%s' created by '%s'\n", feed.Name, feed.Url, user.Name)
	return nil
}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to retrieve feeds: %v", err)
	}
	for _, feed := range feeds {
		username, err := s.db.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Failed to retrieve username with ID: %v", err)
		}
		fmt.Printf("> '%s' @ '%s' by '%s'\n", feed.Name, feed.Url, username)
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
