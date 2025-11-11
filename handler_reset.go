package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.ResetUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to reset users table: %w", err)
	}
	fmt.Println("users table was reset")
	return nil
}
