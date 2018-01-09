package session

import (
	"errors"

	"github.com/utrack/woofer/domain"
)

// Storage stores info about current user logins (pairing session and user ids together).
type Storage interface {
	// IDForSession retrieves user ID for given session ID.
	// It renews any TTLs for the record if it was found.
	IDForSession(sessID string) (domain.UserID, error)
	// SaveID creates a new session and assigns userID to it.
	// It generates a unique session ID internally.
	SaveID(uid domain.UserID) (string, error)
	Delete(sessID string) error
}

// ErrNotFound is returned by Storage if session was not found by sessID requested.
var ErrNotFound = errors.New("session not found by key")
