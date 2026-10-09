package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/sharipovr/blog-aggregator-go/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login expects a single argument, the username")
	}
	name := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return err // No such user or DB problem
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("user has been set / logged in")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("register expects a single argument, the username")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0]})
	if err != nil {
		return err
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("user has been registered / created %+v\n", user)
	return nil
}

func handlerReset(s *state, _ command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("All users have been deleted. Reset complete.")
	return nil
}
