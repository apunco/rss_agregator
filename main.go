package main

import (
	"fmt"
	"os"

	config "github.com/apunco/go/rss_agregator/internal/config"
)

func main() {
	conf, err := config.Read()

	if err != nil {
		fmt.Println(err)
		return
	}

	appState := state{
		cfg: &conf,
	}

	cmds := commands{
		commands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	args := os.Args

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: command required")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&appState, cmd)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
}
