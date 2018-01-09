package sqlite

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utrack/woofer/domain"
	"github.com/utrack/woofer/service/internal/storage"
)

type subsStorage struct {
	c *conn
}

var _ storage.SubsStorage = &subsStorage{}

func (s *subsStorage) Subscribe(ctx context.Context, from, to domain.UserID) error {
	_, err := s.c.ExecContext(ctx,
		`INSERT INTO subs (sfrom,sto) VALUES (?,?)`, from, to)
	return errors.Wrap(err, "error returned from sqlite")
}

func (s *subsStorage) Unsubscribe(ctx context.Context, from, to domain.UserID) error {
	_, err := s.c.ExecContext(ctx,
		`DELETE FROM subs WHERE (sfrom,sto) = (?,?)`, from, to)
	return errors.Wrap(err, "error returned from sqlite")
}

func (s *subsStorage) Subs(ctx context.Context, user domain.UserID) ([]domain.UserID, error) {
	rows, err := s.c.sq.QueryContext(ctx, `SELECT sto FROM subs WHERE sfrom = ?`, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := make([]domain.UserID, 0, 100)
	for rows.Next() {
		var user domain.UserID
		err := rows.Scan(&user)
		if err != nil {
			return nil, errors.Wrap(err, "error when scanning rows from SQL")
		}
		ret = append(ret, user)
	}

	return ret, nil
}

func (s *subsStorage) Subbed(ctx context.Context, user domain.UserID) ([]domain.UserID, error) {
	rows, err := s.c.sq.QueryContext(ctx, `SELECT sfrom FROM subs WHERE sto = ?`, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := make([]domain.UserID, 0, 100)
	for rows.Next() {
		var user domain.UserID
		err := rows.Scan(&user)
		if err != nil {
			return nil, errors.Wrap(err, "error when scanning rows from SQL")
		}
		ret = append(ret, user)
	}

	return ret, nil
}
