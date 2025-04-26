package db

import (
	"encoding/json"

	"github.com/shahin-salehi/website/internal/types"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Crud interface {
	// User Auth
	CreateUser(newUser types.NewUser) (int64, error)
	GetUserByEmail(email string) (*types.User, error)
	GetUserByID(id int64) (*types.UserMeta, error)
	DeleteUser(userId int64) error

	// Messages
	CreateMessage(newMsg types.NewMessage) (string, error)
	GetUserMessages(userId int64) ([]types.Message, error)
	DeleteMessage(userId int64, messageId string) error

	// GDPR
	ExportUserData(userId int64) (*json.RawMessage, error)
	DeleteUserData(userId int64) error
}

type crud struct {
	db *pgxpool.Pool
}

func NewRepo(conn *pgxpool.Pool) *crud {
	return &crud{db: conn}
}
