package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Atm-0s/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func argCheck(
	cmd command,
	tooFew string,
) error {
	if len(cmd.Args) == 0 {
		return errors.New(tooFew)
	}
	if len(cmd.Args) > 1 {
		return errors.New("too many arguments provided")
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	err := argCheck(
		cmd,
		"login requires a username",
	)
	if err != nil {
		return err
	}

	username := cmd.Args[0]
	ctx := context.Background()
	_, err = s.db.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("user %s does not exist", username)
	}

	err = s.cfgPtr.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("User has been set to %s\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	err := argCheck(
		cmd,
		"register requires a name",
	)
	if err != nil {
		return err
	}
	ctx := context.Background()
	name := cmd.Args[0]
	_, err = s.db.GetUser(ctx, name)
	if err == nil {
		return fmt.Errorf("user %s already exists", name)
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	newUser, err := s.db.CreateUser(ctx, userParams)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	err = handlerLogin(s, cmd)
	if err != nil {
		return err
	}
	log.Printf("user was created: id=%s name=%s created_at=%s", newUser.ID, newUser.Name, newUser.CreatedAt)
	return nil
}
