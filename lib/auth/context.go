package auth

import (
	"context"

	"github.com/utrack/woofer/domain"
	"github.com/utrack/woofer/lib/bizerr"
)

const ctxUserIDKey = `auth.userID`

// ErrNoLogin returned if user is not logged in.
var ErrNoLogin = bizerr.New("No login info for the request", bizerr.ErrorUnauthorized)

// UserID returns user ID of a user making the request.
func UserID(ctx context.Context) (domain.UserID, error) {
	v := ctx.Value(ctxUserIDKey)
	if v == nil {
		return 0, ErrNoLogin
	}
	return v.(domain.UserID), nil
}

// SetUserID sets user ID for this request.
func SetUserID(ctx context.Context, uid domain.UserID) context.Context {
	return context.WithValue(ctx, ctxUserIDKey, uid)
}
