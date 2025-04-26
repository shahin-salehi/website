package db

import (
	"context"
	"encoding/json"
	"log/slog"
)

// ExportUserData returns JSON-formatted user data (for GDPR export)
func (c *crud) ExportUserData(userId int64) (*json.RawMessage, error) {
	sql_stmt := `SELECT export_user_data($1);`

	var json json.RawMessage
	err := c.db.QueryRow(context.Background(), sql_stmt, userId).Scan(&json)
	if err != nil {
		slog.Error("failed to export user data", slog.Any("function", "ExportUserData"), slog.Any("error", err))
		return nil, err
	}

	return &json, nil
}
