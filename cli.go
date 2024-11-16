package main

import (
	"context"
	"fmt"
	"time"

	"github.com/0x4D5352/bootdev_blog_aggregator/internal/config"
	"github.com/0x4D5352/bootdev_blog_aggregator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("username is required to login!")
	}
	// TODO: sanitize this input!
	username := cmd.arguments[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		if user.Name == "" {
			return fmt.Errorf("user does not exist!")
		}
		return err
	}
	err = s.config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("%s set as user!\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("name is required to register!")
	}
	t := time.Now().UTC()
	name := cmd.arguments[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: t,
		UpdatedAt: t,
		Name:      name,
	})
	if err != nil {
		return err
	}
	s.config.SetUser(user.Name)
	fmt.Printf("Created user %s!\nUser Details:\n%+v\n", user.Name, user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	// TODO: decide on if you want to mask any of the error codes or add more safety rails
	return s.db.ResetUsers(context.Background())
}

func handlerGetUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("All Users:")
	for _, user := range users {
		name := user.Name
		if name == s.config.CurrentUserName {
			name = fmt.Sprintf("%s (current)", name)
		}
		fmt.Printf("* %s\n", name)
	}
	return nil
}

func handlerAgg(s *state, _ command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}
	fmt.Printf("feed:\n%+v\n", feed)
	return nil
}

func hanlderAddFeed(s *state, cmd command) error {
	name := cmd.arguments[0]
	url := cmd.arguments[1]
	return nil
}

type commands struct {
	maps map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	if _, ok := c.maps[name]; ok {
		fmt.Printf("%s already registered!\n", name)
		return
	}
	c.maps[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if _, ok := c.maps[cmd.name]; !ok {
		return fmt.Errorf("%s not a registered command!", cmd.name)
	}
	return c.maps[cmd.name](s, cmd)
}
