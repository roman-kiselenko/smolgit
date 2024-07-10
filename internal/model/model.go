package model

import "time"

type User struct {
	ID        *int64    `db:"id"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

type Key struct {
	ID        *int64    `db:"id"`
	UserID    string    `db:"user_id"`
	KeyID     string    `db:"key_id"`
	CreatedAt time.Time `db:"created_at"`
}

type Repository struct {
	ID        *int64    `db:"id"`
	UserID    int64     `db:"user_id"`
	Path      string    `db:"path"`
	CreatedAt time.Time `db:"created_at"`
}
