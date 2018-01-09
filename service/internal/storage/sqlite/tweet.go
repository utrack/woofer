package sqlite

import (
	"context"

	sqlite "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/utrack/woofer/domain"
	"github.com/utrack/woofer/lib/bizerr"
	"github.com/utrack/woofer/service/internal/storage"
)

type tweetStorage struct {
	c *conn
}

var _ storage.TweetStorage = &tweetStorage{}

func (ts *tweetStorage) Tweet(ctx context.Context, t domain.Tweet) (uint64, error) {
	res, err := ts.c.ExecContext(ctx,
		`INSERT INTO tweets (uid,created_at,text) VALUES (?,?,?)`, t.From, t.At, t.Text)
	if err != nil {
		return 0, errors.Wrap(err, "error returned from sqlite")
	}

	ret, _ := res.LastInsertId()
	return uint64(ret), nil
}

func (ts *tweetStorage) ByID(ctx context.Context, id uint64) (domain.Tweet, error) {
	var ret domain.Tweet
	row := ts.c.sq.QueryRowContext(ctx, `SELECT id,uid,created_at,text FROM tweets WHERE id = ?`, id)
	err := row.Scan(&ret.ID, &ret.From, &ret.At, &ret.Text)
	if err != nil {
		if err, ok := err.(sqlite.Error); ok && err.Code == sqlite.ErrEmpty {
			return ret, bizerr.New("tweet was not found", bizerr.ErrorNotFound)
		}
	}
	return ret, errors.Wrap(err, "error when scanning tweet")

}

func (ts *tweetStorage) GetPageForUser(ctx context.Context, user domain.UserID, fromTweetID uint64, len uint) ([]domain.TweetWithUsername, error) {
	rows, err := ts.c.sq.QueryContext(ctx, `
SELECT t.id,u.nickname,created_at,text
FROM tweets t
JOIN users u
 ON t.uid = u.id
WHERE
uid IN
 (SELECT sto FROM subs WHERE sfrom = ?)
AND t.id > ?
ORDER BY t.id
LIMIT ?`, user, fromTweetID, len)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := make([]domain.TweetWithUsername, 0, len)

	for rows.Next() {
		var tweet domain.TweetWithUsername
		err := rows.Scan(&tweet.ID, &tweet.From, &tweet.At, &tweet.Text)
		if err != nil {
			return nil, errors.Wrap(err, "error when scanning rows from SQL")
		}
		ret = append(ret, tweet)
	}
	return ret, nil
}

func (ts *tweetStorage) GetPageForProfile(ctx context.Context, user domain.UserID, fromTweetID uint64, len uint) ([]domain.TweetWithUsername, error) {
	rows, err := ts.c.sq.QueryContext(ctx, `
SELECT t.id,u.nickname,created_at,text
FROM tweets t
JOIN users u
 ON t.uid = u.id
WHERE
uid = ?
AND t.id > ?
ORDER BY t.id
LIMIT ?`, user, fromTweetID, len)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := make([]domain.TweetWithUsername, 0, len)

	for rows.Next() {
		var tweet domain.TweetWithUsername
		err := rows.Scan(&tweet.ID, &tweet.From, &tweet.At, &tweet.Text)
		if err != nil {
			return nil, errors.Wrap(err, "error when scanning rows from SQL")
		}
		ret = append(ret, tweet)
	}
	return ret, nil
}
