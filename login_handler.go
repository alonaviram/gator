package main

import (
	"context"
	"errors"
)

func handlerLogin(s *state, cmd command) error {
	args := cmd.args
	if len(args) != 1 {
		return errors.New("need username as argument")
	}
	name := args[0]
	_, e := s.db.GetUser(context.Background(), name)
	if e != nil {
		return e
	}
	err := s.config.SetUser(name)
	if err != nil {
		return err
	}
	return nil
}
