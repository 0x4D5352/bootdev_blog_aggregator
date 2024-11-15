package login

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

func handlerLogin(s *state, cmd command) error {
	if cmd.arguments == nil {
		return fmt.Errorf("expected username; received nil args")
	}
	// TODO: sanitize this input!
	err := s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("%s set as user!", cmd.arguments[0])
	return nil
}
