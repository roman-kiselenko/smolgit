package model

type User struct {
	Name        string
	Permissions string
	Repos       []*Repository
	Keys        []string
}

type Repository struct {
	User *User
	Path string
	Refs []string
	Tags []string
}

func (r Repository) GetFullPath() string {
	return r.User.Name + "/" + r.Path
}
