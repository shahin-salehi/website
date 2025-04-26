package session

import (
	"encoding/base64"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

type Session struct {
	store *sessions.CookieStore
}

const sessionName = "shahin-session"

func InitStore() (*Session, error) {
	secret := os.Getenv("SESSION_KEY")
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		slog.Error("Failed to decode SESSION_KEY", slog.Any("error", err))
		return nil, err

	}
	store := sessions.NewCookieStore(key)
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true, // change to false if testing without HTTPS
	}

	//Ok
	return &Session{store: store}, nil
}

// get our initialized session and store the id of the user
func (s *Session) SetUserID(w http.ResponseWriter, r *http.Request, userID int64) error {
	session, err := s.store.Get(r, sessionName)
	if err != nil {
		slog.Warn("SetUserID returned error and new session, could be decoding error", slog.Any("error", err))
	}

	// set current user in browser session
	session.Values["user_id"] = userID

	// if save to browser fails it will return error
	return session.Save(r, w)
}

// get user id, we'll use this when checking session
// ex validate person before calling db with user sensitive request
func (s *Session) GetUserID(r *http.Request) (int64, bool) {
	session, err := s.store.Get(r, sessionName)
	if err != nil {
		slog.Warn("GetUserID returned error and new session, could be decoding error", slog.Any("error", err))
	}
	id, ok := session.Values["user_id"].(int64)
	return id, ok
}

// clear session essentially during logout or timout (maybe)
func (s *Session) Clear(w http.ResponseWriter, r *http.Request) error {
	session, err := s.store.Get(r, sessionName)
	if err != nil {
		slog.Warn("Clear returned error and new session, could be decoding error", slog.Any("error", err))
	}
	// https://stackoverflow.com/questions/62552103/how-to-remove-a-user-session-from-store-in-gorilla
	session.Options.MaxAge = -1

	// we need to handle this in caller
	return session.Save(r, w)
}
