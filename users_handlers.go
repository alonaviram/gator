package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alonaviram/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	args := cmd.args
	if len(args) != 1 {
		return errors.New("need username as argument")
	}
	name := args[0]
	_, e := s.db.GetUserByName(context.Background(), name)
	if e != nil {
		return e
	}
	err := s.config.SetUser(name)
	if err != nil {
		return err
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	args := cmd.args
	if len(args) != 1 {
		return errors.New("need username as argument")
	}
	name := cmd.args[0]
	_, err := s.db.GetUserByName(context.Background(), name)
	if err == nil {
		return errors.New("user already exists")
	}

	u, e := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
		},
	)
	if e != nil {
		return e
	}
	s.config.SetUser(u.Name)
	fmt.Printf("created user: %v", u)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	currentUser := s.config.CurrentUserName
	for _, u := range users {
		var post string
		if currentUser == u.Name {
			post = " (current)"
		}
		fmt.Printf("%v%v\n", u.Name, post)
	}

	return nil
}
