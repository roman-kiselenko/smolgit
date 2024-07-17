package model

type User struct {
	Login    string
	Password string
	Repos    []*Repository
	Keys     []string
}

type Repository struct {
	User *User
	Path string
	Refs []string
	Tags []string
}
