package main

import (
	config "github.com/apunco/go/rss_agregator/internal/config"
	"github.com/apunco/go/rss_agregator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}
