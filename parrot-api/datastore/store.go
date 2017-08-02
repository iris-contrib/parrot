// Package datastore specifies and implements the datastore required by the API.
package datastore

import (
	"database/sql"

	dbErrors "github.com/iris-contrib/parrot/parrot-api/datastore/errors"
	"github.com/iris-contrib/parrot/parrot-api/datastore/postgres"
	"github.com/iris-contrib/parrot/parrot-api/model"
)

// Store is the interface that datastores must implement.
type Store interface {
	model.LocaleStorer
	model.ProjectStorer
	model.ProjectLocaleStorer
	model.UserStorer
	model.ProjectUserStorer
	model.ProjectClientStorer
	Ping() error
	Close() error
	MigrateUp(string) error
	MigrateDown(string) error
}

// Datastore is the provided Store implementation.
type Datastore struct {
	Store
}

// NewDatastore creates and configures a new datastore based on the
// parameter name and the connection url.
// Currently only 'postgres' is supported.
func NewDatastore(name string, url string) (*Datastore, error) {
	var ds *Datastore

	switch name {
	case "postgres":
		conn, err := sql.Open("postgres", url)
		if err != nil {
			return nil, err
		}

		p := &postgres.PostgresDB{DB: conn}
		// TODO: debug refused connections when db connections > 1
		p.SetMaxIdleConns(1)
		p.SetMaxOpenConns(1)

		ds = &Datastore{p}
	default:
		return nil, dbErrors.ErrNotImplemented
	}

	return ds, nil
}
