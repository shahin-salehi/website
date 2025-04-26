package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/shahin-salehi/website/internal/types"

	"golang.org/x/crypto/bcrypt"
)

// add client info on expected failures
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	// check
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// unmarshal
	payload := new(types.Login)
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		slog.Error("failed to unmarshal login payload", slog.Any("error", err))
		writeJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// for ease
	email := payload.Email
	password := payload.PasswordHash

	if email == "" || password == "" {
		writeJSONError(w, "Missing email or password", http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUserByEmail(email)
	if err != nil {
		writeJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		writeJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = h.session.SetUserID(w, r, user.Id)
	if err != nil {
		slog.Error("Failed to store session", slog.Any("error", err))
		writeJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"redirect": "/",
	})

}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	// check
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.session.Clear(w, r)
	if err != nil {
		slog.Error("failed to clear session", slog.Any("error", err))
		writeJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	// check
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// unmarshal
	payload := new(types.NewUser)
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		slog.Error("failed to unmarshal signup payload", slog.Any("error", err))
		writeJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// for ease
	username := payload.Username
	email := payload.Email
	password := payload.PasswordHash

	if email == "" || password == "" || username == "" {
		writeJSONError(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", slog.Any("error", err))
		writeJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	user := types.NewUser{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}

	_, err = h.db.CreateUser(user)
	if err != nil {
		slog.Error("db returned error when creating user", slog.Any("error", err))
		writeJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"server": "Account created, please login.",
	})
}

func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := h.session.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.db.DeleteUser(userID)
	if err != nil {
		slog.Error("Failed to delete user", slog.Any("userID", userID), slog.Any("error", err))
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	_ = h.session.Clear(w, r)

	http.Redirect(w, r, "/account-deleted", http.StatusSeeOther)
}

func (h *Handler) GetUserData(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := h.session.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// get data
	data, err := h.db.ExportUserData(userID)
	if err != nil {
		slog.Error("failed to get user data", slog.Any("userID", userID), slog.Any("error", err))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// timestamp
	now := time.Now().UTC()
	filename := "mydata-" + strings.ReplaceAll(now.Format(time.DateTime), " ", "_") + ".json"

	w.Header().Set("Content-Type", "application/json")
	// download
	w.Header().Set("Content-Disposition", `attachment; filename=`+filename+`"`)
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("failed to encode user data", slog.Any("userID", userID))
	}
}

func (h *Handler) DeleteUserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := h.session.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.db.DeleteUserData(userID)
	if err != nil {
		slog.Error("Failed to delete user data", slog.Any("userID", userID), slog.Any("error", err))
		// the client call doesn't execpt json so we get away here with not
		// responding with proper 500
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/data-deleted", http.StatusSeeOther)
}
