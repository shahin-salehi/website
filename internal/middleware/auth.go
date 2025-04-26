package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/shahin-salehi/website/internal/db"
	"github.com/shahin-salehi/website/internal/session"
)

// maybe merge with handler

type Middlware struct {
	db db.Crud
}

func NewMiddleware(repo db.Crud) *Middlware {
	return &Middlware{db: repo}
}

// use as middleware
func (m *Middlware) RequireAuth(sess *session.Session, next http.Handler) http.Handler {
	// wrap our function to satsify http.Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check session
		_, ok := sess.GetUserID(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		// cont
		next.ServeHTTP(w, r)
	})
}

// what we do here is we check if the user is signed in
// if so, we pass their info
// keep in mind this will call the db everytime they navigate
// the page and given that they can never change their info,
// yet, that's probably not necessary.

// https://github.com/a-h/templ/issues/253

type contextKey string

const ContextUser = contextKey("user")

func (m *Middlware) InjectUser(sess *session.Session, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if session is active
		id, ok := sess.GetUserID(r)
		if !ok {
			// skip, cont
			next.ServeHTTP(w, r)
			return
		}

		// get metadata to inject
		user, err := m.db.GetUserByID(id)
		if err != nil {
			// this shouldn't happen
			slog.Error("session returned id but user doesn't exists in database, this shouldn't happen", slog.Any("error", err))

			// cont chain
			next.ServeHTTP(w, r)
			return
		}

		// set context
		ctx := context.WithValue(r.Context(), ContextUser, user.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middlware) RedirectAuth(sess *session.Session, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if session is active don't continue send them out
		_, ok := sess.GetUserID(r)
		if ok {
			// no need for this
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// if no session cont
		next.ServeHTTP(w, r)
	})
}
