package db

import (
	"context"
	"log/slog"

	"github.com/shahin-salehi/website/internal/types"
)

// CreateUser inserts a new user into the database
func (c *crud) CreateUser(newUser types.NewUser) (int64, error) {
	sql_stmt := `SELECT create_user($1, $2, $3);`

	var userID int64
	err := c.db.QueryRow(context.Background(), sql_stmt, newUser.Username, newUser.Email, newUser.PasswordHash).Scan(&userID)
	if err != nil {
		slog.Error("failed to create user", slog.Any("function", "CreateUser"), slog.Any("error", err))
		return 0, err
	}

	return userID, nil
}

// GetUserByEmail fetches the user by email for login
func (c *crud) GetUserByEmail(email string) (*types.User, error) {
	sql_stmt := `SELECT * FROM get_user_by_email($1);`

	row := c.db.QueryRow(context.Background(), sql_stmt, email)

	user := new(types.User)
	err := row.Scan(&user.Id, &user.PasswordHash)
	if err != nil {
		slog.Error("failed to fetch user", slog.Any("function", "GetUserByEmail"), slog.Any("email", email), slog.Any("error", err))
		return nil, err
	}

	user.Email = email // Add it manually because function doesn't return it
	return user, nil
}

// GetUserByEmail fetches the user by email for login
func (c *crud) GetUserByID(id int64) (*types.UserMeta, error) {
	sql_stmt := `SELECT * FROM get_user_by_id($1);`

	row := c.db.QueryRow(context.Background(), sql_stmt, id)

	user := new(types.UserMeta)
	err := row.Scan(&user.Id, &user.Username)
	if err != nil {
		slog.Error("failed to fetch user meta", slog.Any("function", "GetUserByID"), slog.Any("error", err))
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes the user and cascades their data
func (c *crud) DeleteUser(userId int64) error {
	sql_stmt := `SELECT delete_user($1);`

	_, err := c.db.Exec(context.Background(), sql_stmt, userId)
	if err != nil {
		slog.Error("failed to delete user", slog.Any("function", "DeleteUser"), slog.Any("error", err))
		return err
	}

	return nil
}

func (c *crud) DeleteUserData(userID int64) error {
	sqlStmt := `SELECT delete_user_messages($1);`

	_, err := c.db.Exec(context.Background(), sqlStmt, userID)
	if err != nil {
		slog.Error("failed to delete user data", slog.Any("function", "DeleteUserData"), slog.Any("error", err))
		return err
	}

	return nil
}
