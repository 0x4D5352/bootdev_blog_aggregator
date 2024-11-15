package main

import (
	"fmt"
	"github.com/0x4D5352/bootdev_blog_aggregator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	maps map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if cmd.arguments == nil {
		return fmt.Errorf("expected username; received nil args")
	}
	// TODO: sanitize this input!
	username := cmd.arguments[0]
	err := s.config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("%s set as user!\n", username)
	return nil
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
