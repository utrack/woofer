/*Package storage provides storage primitives for the service.
 */
package storage

import (
	"context"

	"github.com/utrack/woofer/domain"
)

type TweetLister interface {
	ByID(context.Context, uint64) (domain.Tweet, error)
	GetPageForProfile(ctx context.Context, user domain.UserID, fromTweetID uint64, len uint) ([]domain.TweetWithUsername, error)
	GetPageForUser(ctx context.Context, user domain.UserID, fromTweetID uint64, len uint) ([]domain.TweetWithUsername, error)
}

type TweetSaver interface {
	Tweet(context.Context, domain.Tweet) (uint64, error)
}

type TweetStorage interface {
	TweetLister
	TweetSaver
}

type UserSaver interface {
	// New saves a user, returning its ID on success.
	New(context.Context, domain.UserWithPassword) (domain.UserID, error)
	// Save modifies an existing user's data.
	Save(context.Context, domain.User) error
}

type UserLister interface {
	GetByIds(context.Context, []domain.UserID) ([]domain.User, error)
	GetByNickname(context.Context, string) (domain.User, error)
}

type UserStorage interface {
	UserLister
	UserSaver
}

type PasswordManager interface {
	PasswordCheck(context.Context, string, string) (bool, error)
}

// SubsLister lists users' subscriptions.
type SubsLister interface {
	// Subs returns a user's subscriptions.
	Subs(context.Context, domain.UserID) ([]domain.UserID, error)
	// Subbed returns all users who's subscribed to a given user.
	Subbed(context.Context, domain.UserID) ([]domain.UserID, error)
}

// SubsSaver stores info about users' subscriptions.
type SubsSaver interface {
	// Subscribe subscribes first user to a second user.
	Subscribe(context.Context, domain.UserID, domain.UserID) error
	Unsubscribe(context.Context, domain.UserID, domain.UserID) error
}

// SubsStorage is a subscription data storage.
type SubsStorage interface {
	SubsLister
	SubsSaver
}
