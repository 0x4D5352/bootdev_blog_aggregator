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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, currentUser)
	}
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

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("duration required to aggregate feeds!")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("name and url required to add feed!")
	}
	name := cmd.arguments[0]
	url := cmd.arguments[1]
	t := time.Now().UTC()
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: t,
		UpdatedAt: t,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("added feed:\n%+v\n", feed)
	err = handlerFollow(s, command{
		name:      "internal follow",
		arguments: []string{url},
	},
		user)
	if err != nil {
		return err
	}
	return nil
}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("all feeds:")
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("feed name: %s, feed URL: %s,added by: %s\n", feed.Name, feed.Url, user.Name)

	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("url required to follow feed!")
	}
	url := cmd.arguments[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}
	t := time.Now().UTC()
	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: t,
		UpdatedAt: t,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s followed %s!", follow.UserName, follow.FeedName)
	return nil
}

func handlerFollowing(s *state, _ command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}
	fmt.Println("currently following:")
	for _, follow := range follows {
		fmt.Printf("- %s\n", follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("url required to unfollow feed!")
	}
	url := cmd.arguments[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}
	err = s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("unfollowed %s!\n", feed.Name)
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
