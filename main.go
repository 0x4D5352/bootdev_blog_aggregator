package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/0x4D5352/bootdev_blog_aggregator/internal/config" // "github.com/0x4D5352/bootdev_blog_aggregator/cmd/login"
	"github.com/0x4D5352/bootdev_blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Read()
	db, err := sql.Open("postgres", cfg.DbUrl)
	dbQueries := database.New(db)
	if err != nil {
		log.Fatal(err)
	}
	s := state{
		db:     dbQueries,
		config: &cfg,
	}
	cmds := commands{
		maps: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("error: not enough arguments!")
		os.Exit(1)
	}
	name := args[1]
	var arguments []string
	if len(args) > 2 {
		arguments = os.Args[2:]
	}
	cmd := command{
		name:      name,
		arguments: arguments,
	}
	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
