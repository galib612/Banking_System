package authdb

import (
	"database/sql"
)

type Store interface {
	Querier
	// todo : add function to this interface
}

//Store provides all functions to execute SQl queries for authentication and Login
type SQLStore struct {
	*Queries // Inherit all the functions thats implemented by Queries.
	db       *sql.DB
}

//NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}
