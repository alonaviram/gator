package main

import (
	"context"

	"github.com/alonaviram/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	a := func(s *state, cmd command) error {
		currentUser, e := s.db.GetUserByName(context.Background(), s.config.CurrentUserName)
		if e != nil {
			return e
		}
		return handler(s, cmd, currentUser)
	}
	return a
}
