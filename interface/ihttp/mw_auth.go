package ihttp

import (
	"errors"
	"net/http"

	"github.com/utrack/woofer/lib/auth"
	"github.com/utrack/woofer/lib/session"
)

const cookieSessID = "sessid"

// UserAuthCtx injects user ID to the context.
func UserAuthCtx(sessStorage session.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If we need to store more data than current user ID -
			// JWT+JWT.id existence check in sessstorage might be more useful.
			// exists check is needed to logout users reliably.

			sessionID, err := r.Cookie(cookieSessID)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			uid, err := sessStorage.IDForSession(sessionID.Value)
			if err == nil {
				r = r.WithContext(auth.SetUserID(r.Context(), uid))
			} else {
				// TODO log problem if err != ErrNotFound
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAuth blocks access to a handler if user is not logged in.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.UserID(r.Context())
		if err != nil {
			renderError(w, errors.New("Login required"), 403)
			return
		}
		next.ServeHTTP(w, r)
	})
}
