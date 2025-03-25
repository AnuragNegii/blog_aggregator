package main

import (
	"github.com/AnuragNegii/blog_aggregator/internal/config"
	"github.com/AnuragNegii/blog_aggregator/internal/database"
)

type state struct{
    db *database.Queries
    config *config.Config
}
