package db

import (
	"context"
	"log/slog"

	"github.com/shahin-salehi/website/internal/types"

	"github.com/jackc/pgx/v5"
)

// CreateMessage inserts a new message for a user
func (c *crud) CreateMessage(newMsg types.NewMessage) (string, error) {
	sql_stmt := `SELECT create_message($1, $2);`

	var id string
	err := c.db.QueryRow(context.Background(), sql_stmt, newMsg.UserId, newMsg.Content).Scan(&id)
	if err != nil {
		slog.Error("failed to create message", slog.Any("function", "CreateMessage"), slog.Any("error", err))
		return "", err
	}

	return id, nil
}

// GetUserMessages returns all messages for a given user
func (c *crud) GetUserMessages(userId int64) ([]types.Message, error) {
	sql_stmt := `SELECT * FROM get_user_messages($1);`

	rows, err := c.db.Query(context.Background(), sql_stmt, userId)
	if err != nil {
		slog.Error("failed to fetch messages", slog.Any("function", "GetUserMessages"), slog.Any("error", err))
		return nil, err
	}

	structured, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.Message])
	if err != nil {
		slog.Error("failed to collect messages", slog.Any("function", "GetUserMessages"), slog.Any("error", err))
		return nil, err
	}

	return structured, nil
}

// DeleteMessage deletes a specific message if user owns it
func (c *crud) DeleteMessage(userId int64, messageId string) error {
	sql_stmt := `SELECT delete_message($1, $2);`

	_, err := c.db.Exec(context.Background(), sql_stmt, userId, messageId)
	if err != nil {
		slog.Error("failed to delete message", slog.Any("function", "DeleteMessage"), slog.Any("error", err))
		return err
	}

	return nil
}
