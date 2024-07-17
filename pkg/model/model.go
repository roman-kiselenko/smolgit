package model

type User struct {
	User  string
	Repos []*Repository
	Keys  []string
}

type Repository struct {
	User *User
	Path string
	Refs []string
	Tags []string
}

func (r Repository) GetFullPath() string {
	return r.User.User + "/" + r.Path
}
