/*Package sqlite provides sqlite3-backed storage.
 */
package sqlite

import (
	"context"
	"database/sql"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/sqlite3"
	_ "github.com/mattes/migrate/source/file"
	"github.com/pkg/errors"
)

const passwordHashCost = 12

// Storage implements complete storage.Storage using sqlite as a backend.
type Storage struct {
	userStorage
	tweetStorage
	subsStorage
}

// New creates a new sqlite-backed storage.
func New(connstring string, migrations string) (*Storage, error) {
	db, err := sqlx.Connect("sqlite3", connstring)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't init sqlite3 connection")
	}

	dri, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't init migration driver")
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrations,
		"sqlite", dri)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create migration engine")
	}
	err = m.Steps(1)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't run migrations")
	}

	c := &conn{sq: db}
	return &Storage{
		userStorage{
			c: c,
		},
		tweetStorage{
			c: c,
		},
		subsStorage{
			c: c,
		},
	}, nil
}

type conn struct {
	sq *sqlx.DB
	// need to lock a connection since sqlite allows single concurrent write
	// op only
	mtx sync.Mutex
}

func (c *conn) ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.sq.ExecContext(ctx, q, args...)
}

func (c *conn) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.sq.GetContext(ctx, dest, query, args...)
}
