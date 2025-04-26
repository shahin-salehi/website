package handler

import (
	"encoding/json"
	"net/http"

	"github.com/shahin-salehi/website/internal/db"
	"github.com/shahin-salehi/website/internal/session"
)

type Handler struct {
	db      db.Crud
	session *session.Session
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// init
func NewHandler(repo db.Crud, sess *session.Session) *Handler {
	return &Handler{db: repo, session: sess}
}

// helper
func writeJSONError(w http.ResponseWriter, msg string, status int) {
	// for server logs

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}
