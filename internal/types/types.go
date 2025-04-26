package types

import "time"

type User struct {
	Id           int64     `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

// user metadata
type UserMeta struct {
	Id       int64  `db:"id"`
	Username string `db:"username"`
}

type NewUser struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}

type Login struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}

type Message struct {
	Id        string    `db:"id"` // UUID
	UserId    int64     `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

type NewMessage struct {
	UserId  int64
	Content string
}

type ContactMessage struct {
	Message string `json:"message"`
}
