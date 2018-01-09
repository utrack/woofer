package service

import (
	"github.com/pkg/errors"
	"github.com/utrack/woofer/service/internal/storage/sqlite"
)

type Config struct {
	SQLiteConnString string
	SQLiteMigrations string
}

// Bootstrap returns a Woofer service.
func Bootstrap(cfg Config) (*Woofer, error) {

	storage, err := sqlite.New(cfg.SQLiteConnString, cfg.SQLiteMigrations)
	if err != nil {
		return nil, errors.Wrap(err, "storage init failed")
	}
	// Normally we'd provide some configuration for the service there
	// but this is a code challenge so
	return &Woofer{
		tweetStorage: storage,
		userStorage:  storage,
		subStorage:   storage,
		passCheck:    storage,
	}, nil
}
