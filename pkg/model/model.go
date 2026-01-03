package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Name        string        `json:"name,omitempty"`
	Permissions string        `json:"permissions,omitempty"`
	Repos       []*Repository `json:"repos,omitempty"`
	Keys        []string      `json:"keys,omitempty"`
}

type Repository struct {
	User *User    `json:"user,omitempty"`
	Path string   `json:"path,omitempty"`
	Refs []string `json:"refs,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (r Repository) GetFullPath() string {
	return r.User.Name + "/" + r.Path
}
