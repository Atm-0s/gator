package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("login requires a username")
	}
	if len(cmd.Args) > 1 {
		return errors.New("too many arguments provided, login expects one username")
	}
	username := cmd.Args[0]

	err := s.cfgPtr.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("User has been set to %s\n", username)
	return nil
}
