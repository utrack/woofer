package inmemsessions

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/utrack/woofer/domain"
	"github.com/utrack/woofer/lib/session"
)

// Storage implements session.Storage.
type Storage struct {
	ttl time.Duration
	c   *cache.Cache
}

// sessionIDLength is a length of generated session ID.
const sessionIDLength = 64

var _ session.Storage = &Storage{}

// New creates new Storage.
func New(ttl time.Duration) *Storage {
	return &Storage{
		ttl: ttl,
		c:   cache.New(ttl, ttl*2),
	}
}

// IDForSession implements session.Storage.
func (s *Storage) IDForSession(sessID string) (domain.UserID, error) {
	got, ok := s.c.Get(sessID)
	if !ok {
		return 0, session.ErrNotFound
	}
	// renew TTL
	ret := got.(domain.UserID)
	s.c.Set(string(sessID), ret, s.ttl)

	return ret, nil
}

// SaveID implements session.Storage.
func (s *Storage) SaveID(uid domain.UserID) (string, error) {
	sessID := randString(sessionIDLength)
	return sessID, errors.Wrap(s.c.Add(string(sessID), uid, s.ttl), "error returned from go-cache")
}

// Delete implements session.Storage.
func (s *Storage) Delete(sessID string) error {
	s.c.Delete(sessID)
	return nil
}
