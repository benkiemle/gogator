package main

import (
	"database/sql"

	"github.com/benkiemle/gogator/internal/config"
	"github.com/benkiemle/gogator/internal/database"
)

type state struct {
	db         *database.Queries
	config     *config.Config
	connection *sql.DB
}

func NewState() (state, error) {
	configuration, err := config.Read()
	if err != nil {
		return state{}, err
	}

	db, err := sql.Open("postgres", configuration.ConnectionString)
	if err != nil {
		return state{}, err
	}

	dbQueries := database.New(db)

	s := state{
		config:     configuration,
		db:         dbQueries,
		connection: db,
	}
	return s, nil
}
