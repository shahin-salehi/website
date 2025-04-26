package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shahin-salehi/website/internal/types"
)

func (h *Handler) Message() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// unmarshal
		payload := new(types.ContactMessage)
		defer r.Body.Close()

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			slog.Error("failed to unmarshal login payload", slog.Any("error", err))
			writeJSONError(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		//ease
		message := payload.Message

		if message == "" {
			writeJSONError(w, "message cannot be empty", http.StatusBadRequest)
			return
		}

		// get user id form session
		userID, ok := h.session.GetUserID(r)
		if !ok {
			slog.Error("session couldn't return userid, this shouldn't happen")
			writeJSONError(w, "something went wrong", http.StatusInternalServerError)
		}
		newMsg := types.NewMessage{
			UserId:  userID,
			Content: message,
		}
		_, err = h.db.CreateMessage(newMsg)
		if err != nil {
			slog.Error("Database failed to create message", slog.Any("error", err))
			return
		}

		//ok
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		resp := `{"server":"message sent!"}`
		w.Write([]byte(resp))

	})
}
