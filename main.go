package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	config "github.com/apunco/go/rss_agregator/internal/config"
	"github.com/apunco/go/rss_agregator/internal/database"
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

	db, err := sql.Open("postgres", appState.cfg.DbUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	appState.db = database.New(db)

	cmds := commands{
		commands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", resetHandler)
	cmds.register("users", getUsersHandler)
	cmds.register("agg", getFeedHandler)

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
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		os.Exit(1)
	}
}
