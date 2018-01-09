/*Package service provide woofer's service (and its business logic).
 */
package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/utrack/woofer/domain"
	"github.com/utrack/woofer/lib/auth"
	"github.com/utrack/woofer/lib/bizerr"
	"github.com/utrack/woofer/service/internal/storage"
)

// Woofer provides this service's business logic.
type Woofer struct {
	tweetStorage storage.TweetStorage
	userStorage  storage.UserStorage
	subStorage   storage.SubsStorage
	passCheck    storage.PasswordManager
}

var (
	ErrIncorrectLogin = bizerr.New("Incorrect login or password", bizerr.ErrorUnauthorized)
)

// Tweet posts a new tweet.
func (w Woofer) Tweet(ctx context.Context, text string) (uint64, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "couldn't get UserID for request")
	}
	var t domain.Tweet
	if len(text) == 0 {
		return 0, bizerr.New("tweet cannot be empty", bizerr.ErrorUserInput)
	}

	t.From = userID
	t.At = time.Now()
	t.Text = text

	ret, err := w.tweetStorage.Tweet(ctx, t)
	return ret, errors.Wrap(err, "couldn't post a tweet")
}

// GetTweetPage returns a page of tweets for the current user.
// User can request next pages giving us a tweet ID to start the page from.
func (w Woofer) GetTweetPage(ctx context.Context, from uint64) ([]domain.TweetWithUsername, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get UserID for request")
	}

	ret, err := w.tweetStorage.GetPageForUser(ctx, userID, from, 30)
	return ret, errors.Wrap(err, "couldn't retrieve tweets")
}

// GetTweetsForProfile returns a tweet list for given user.
func (w Woofer) GetTweetsForProfile(ctx context.Context, user string, from uint64) ([]domain.TweetWithUsername, error) {

	tgt, err := w.userStorage.GetByNickname(ctx, user)
	if err != nil {
		return nil, err
	}

	ret, err := w.tweetStorage.GetPageForProfile(ctx, tgt.ID, from, 30)
	return ret, errors.Wrap(err, "couldn't retrieve tweets")

}

// Subscribe subscribes current user to another one.
func (w Woofer) Subscribe(ctx context.Context, targetNickname string) error {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return errors.Wrap(err, "couldn't get UserID for request")
	}
	tgt, err := w.userStorage.GetByNickname(ctx, targetNickname)
	if err != nil {
		return err
	}
	return w.subStorage.Subscribe(ctx, userID, tgt.ID)
}

// Unsubscribe unsubscribes current user from some other user.
func (w Woofer) Unsubscribe(ctx context.Context, targetNickname string) error {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return errors.Wrap(err, "couldn't get UserID for request")
	}
	tgt, err := w.userStorage.GetByNickname(ctx, targetNickname)
	if err != nil {
		return err
	}
	return w.subStorage.Unsubscribe(ctx, userID, tgt.ID)
}

// Subscriptions returns all user IDs to which a current user is subscribed to.
func (w Woofer) Subscriptions(ctx context.Context) ([]domain.User, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get UserID for request")
	}

	subs, err := w.subStorage.Subs(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't retrieve subscriptions list")
	}
	return w.userStorage.GetByIds(ctx, subs)
}

// Subscribers returns all user IDs subscribed to current user.
func (w Woofer) Subscribers(ctx context.Context) ([]domain.User, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get UserID for request")
	}

	subs, err := w.subStorage.Subbed(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't retrieve subscriptions list")
	}
	return w.userStorage.GetByIds(ctx, subs)
}

// UserCreate creates new user.
// Caller should NOT be logged in.
func (w Woofer) UserCreate(ctx context.Context, u domain.UserWithPassword) (domain.UserID, error) {
	_, err := auth.UserID(ctx)
	if err == nil {
		return 0, errors.New("user already logged in")
	}

	if len(u.Password) < 6 {
		return 0, bizerr.New("password can't be shorter than 6 chars", bizerr.ErrorUserInput)
	}

	// TODO more sanity checks etc
	ret, err := w.userStorage.New(ctx, u)
	return ret, errors.Wrap(err, "couldn't save new user")
}

// UserModify modifies an existing user.
func (w Woofer) UserModify(ctx context.Context, u domain.User) error {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return errors.Wrap(err, "couldn't get UserID for request")
	}
	u.ID = userID

	// TODO sanity checks etc
	err = w.userStorage.Save(ctx, u)
	return errors.Wrap(err, "couldn't modify existing user")
}

// UserByNickname returns a user by its nickname.
// Does not require for the caller to be logged in.
func (w Woofer) UserByNickname(ctx context.Context, id string) (domain.User, error) {
	ret, err := w.userStorage.GetByNickname(ctx, id)
	return ret, errors.Wrap(err, "couldn't access user storage")
}

// CheckPassword checks if password matches the username.
func (w Woofer) CheckPassword(ctx context.Context, username string, pass string) error {
	_, err := auth.UserID(ctx)
	if err == nil {
		return bizerr.New("should be logged out to perform this", bizerr.ErrorUnauthorized)
	}
	ok, err := w.passCheck.PasswordCheck(ctx, username, pass)
	if err != nil {
		return err
	}
	if !ok {
		return ErrIncorrectLogin
	}
	return nil
}
