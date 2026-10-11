package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/sharipovr/blog-aggregator-go/internal/config"
	"github.com/sharipovr/blog-aggregator-go/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {

	// Initial read
	var s state
	c, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	s.config = &c

	// Open a connection to the database
	dbURL := s.config.DbUrl
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	// Saving to state
	dbQueries := database.New(db)
	s.db = dbQueries

	// commands (p7 from CH1. l-03)
	var cmds commands
	cmds.handlers = make(map[string]func(*state, command) error)
	cmds.register("register", handlerRegister)
	cmds.register("login", handlerLogin)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Error: not enough arguments provided")
	}

	cmd := command{name: args[1], args: args[2:]}
	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
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
