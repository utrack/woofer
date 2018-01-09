package sqlite

import (
	"context"
	"fmt"
	"strings"

	sqlite "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/utrack/woofer/domain"
	"github.com/utrack/woofer/lib/bizerr"
	"github.com/utrack/woofer/service/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

// userStorage implements UserStorage and PasswordManager.
type userStorage struct {
	c *conn
}

var _ storage.UserStorage = &userStorage{}
var _ storage.PasswordManager = &userStorage{}

func (us *userStorage) New(ctx context.Context, u domain.UserWithPassword) (domain.UserID, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), passwordHashCost)
	if err != nil {
		return 0, errors.Wrap(err, "couldn't generate a hash")
	}

	res, err := us.c.ExecContext(ctx,
		`INSERT INTO users (name,nickname,password) VALUES (?,?,?)`, u.RealName, u.Nickname, pass)
	if err != nil {
		if err, ok := err.(sqlite.Error); ok && err.Code == sqlite.ErrConstraint {
			return 0, bizerr.New("this nickname is already taken", bizerr.ErrorConflict)
		}
		return 0, errors.Wrap(err, "error returned from sqlite")
	}

	ret, _ := res.LastInsertId()
	return domain.UserID(ret), nil
}

func (us *userStorage) Save(ctx context.Context, u domain.User) error {
	_, err := us.c.ExecContext(ctx,
		`UPDATE users SET (name) = (?) WHERE id = ?`, u.RealName, u.ID)
	return errors.Wrap(err, "error returned from sqlite")
}

func (us *userStorage) PasswordCheck(ctx context.Context, nickname string, pass string) (bool, error) {
	var pwdHash []byte
	err := us.c.GetContext(ctx, &pwdHash, `SELECT password FROM users WHERE nickname = ?`, nickname)
	if err != nil {
		if err, ok := err.(sqlite.Error); ok && err.Code == sqlite.ErrEmpty {
			return false, nil
		}
		return false, errors.Wrap(err, "error returned from sqlite")
	}
	err = bcrypt.CompareHashAndPassword(pwdHash, []byte(pass))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "error comparing hashes")
	}
	return true, nil
}

func (us *userStorage) GetByNickname(ctx context.Context, n string) (domain.User, error) {
	var ret domain.User
	row := us.c.sq.QueryRowContext(ctx, `SELECT id,name,nickname FROM users WHERE nickname = ?`, n)
	err := row.Scan(&ret.ID, &ret.RealName, &ret.Nickname)
	if err != nil {
		if err, ok := err.(sqlite.Error); ok && err.Code == sqlite.ErrEmpty {
			return ret, bizerr.New("user was not found", bizerr.ErrorNotFound)
		}
	}
	return ret, errors.Wrap(err, "error when scanning user")
}

func (us *userStorage) GetByIds(ctx context.Context, ids []domain.UserID) ([]domain.User, error) {
	idsText := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
	rows, err := us.c.sq.QueryContext(ctx, `SELECT id,name,nickname FROM users WHERE id IN (?)`, idsText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := make([]domain.User, 0, len(ids))
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.RealName, &user.Nickname)
		if err != nil {
			return nil, errors.Wrap(err, "error when scanning rows from SQL")
		}
		ret = append(ret, user)
	}

	return ret, nil
}
