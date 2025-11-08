package main

import (
	"fmt"
	"os"

	"github.com/Atm-0s/BlogAggregator/internal/config"
)

type state struct {
	cfgPtr *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	s := &state{
		cfgPtr: &cfg,
	}
	c := commands{
		cmdMap: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)

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
