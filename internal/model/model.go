package model

import "time"

type User struct {
	ID        *int64        `db:"id"`
	Login     string        `db:"login"`
	Password  string        `db:"password"`
	Repos     []*Repository `db:"-"`
	CreatedAt time.Time     `db:"created_at"`
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
	User      *User     `db:"-"`
	Path      string    `db:"path"`
	Refs      []string  `db:"-"`
	Tags      []string  `db:"-"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
