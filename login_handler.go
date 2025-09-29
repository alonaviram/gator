package main

import (
	"errors"
)

func handlerLogin(s *state, cmd command) error {
	args := cmd.args
	if len(args) != 1 {
		return errors.New("need username as argument")
	}
	err := s.config.SetUser(args[0])
	if err != nil {
		return err
	}
	return nil
}
