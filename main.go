package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Atm-0s/BlogAggregator/internal/config"
	"github.com/Atm-0s/BlogAggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	cfgPtr *config.Config
}

func main() {
	// Read config
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Open database
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialise the database queries
	dbQueries := database.New(db)

	// Build state with config and queries
	s := &state{
		db:     dbQueries,
		cfgPtr: &cfg,
	}

	// Register commands
	c := commands{
		cmdMap: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerGetUsers)

	// Parse input and run
	input := os.Args
	if len(input) < 2 {
		fmt.Println("no arguments entered")
		os.Exit(1)
	}

	var inputCMD command
	inputCMD.Name = input[1]
	inputCMD.Args = input[2:]

	err = c.run(s, inputCMD)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
