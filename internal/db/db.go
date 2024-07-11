package db

import (
	_ "embed"
	"fmt"
	"log/slog"

	"smolgit/internal/model"
	"smolgit/internal/password"

	"github.com/jmoiron/sqlx"
	//nolint:revive
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v3"
)

type Database struct {
	*sqlx.DB
	logger    *slog.Logger
	rootLogin string
}

//go:embed schema.sql
var schema string

func New(logger *slog.Logger, clictx *cli.Context) (*Database, error) {
	sqliteDatabase, err := sqlx.Connect("sqlite3", clictx.String("db_path"))
	if err != nil {
		return &Database{}, err
	}
	db := &Database{
		sqliteDatabase, logger, clictx.String("root_login"),
	}
	if _, err := db.MustExec(schema).RowsAffected(); err != nil {
		return &Database{}, fmt.Errorf("migrate database scheme %w", err)
	}
	pass, _ := password.HashPassword(clictx.String("root_password"))
	if _, err := db.NamedExec(`INSERT INTO users (login, password) VALUES (:login, :password)`,
		map[string]interface{}{
			"login":    clictx.String("root_login"),
			"password": pass,
		}); err == nil {
		logger.Info("root user created", "login", clictx.String("root_login"), "password", "*******")
	}

	return db, nil
}

// TODO write to tx
func (db *Database) InsertUser(login, pswd, key string) error {
	pass, _ := password.HashPassword(pswd)
	result, err := db.NamedExec(`INSERT INTO users (login, password) VALUES (:login, :password)`,
		map[string]interface{}{
			"login":    login,
			"password": pass,
		})
	if err != nil {
		return err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	return db.InsertUserKey(userID, key)
}

func (db *Database) InsertUserKey(userID int64, key string) error {
	if _, err := db.NamedExec(`INSERT INTO keys (user_id, key_id) VALUES (:user_id, :key_id)`,
		map[string]interface{}{
			"user_id": userID,
			"key_id":  key,
		}); err != nil {
		return err
	}
	return nil
}

func (db *Database) InsertRepo(userID int64, path string) error {
	if _, err := db.NamedExec(`INSERT INTO repositories (user_id, path) VALUES (:user_id, :path)`,
		map[string]interface{}{
			"user_id": userID,
			"path":    path,
		}); err != nil {
		return err
	}
	return nil
}

func (db *Database) RepoExist(userID int64, path string) bool {
	var exist int
	if err := db.Get(&exist, `SELECT 1 FROM repositories WHERE user_id = $1 AND path = $2;`, userID, path); err != nil {
		return false
	}
	return exist == 1
}

func (db *Database) FindUserFromKey(key string) (model.User, error) {
	var user model.User
	if err := db.Get(&user, `SELECT users.id, login FROM users LEFT JOIN keys ON users.id = keys.user_id WHERE keys.key_id = $1;`, key); err != nil {
		return user, err
	}
	return user, nil
}

func (db *Database) FindRootUser() (model.User, error) {
	var user model.User
	if err := db.Get(&user, `SELECT * FROM users WHERE login = $1;`, db.rootLogin); err != nil {
		return user, err
	}
	return user, nil
}

// TODO fix query
func (db *Database) ListUsers() ([]*model.User, error) {
	var users []*model.User
	if err := db.Select(&users, `SELECT * FROM users ORDER BY created_at ASC;`); err != nil {
		return users, err
	}
	userIDs := []int64{}
	for _, u := range users {
		userIDs = append(userIDs, *u.ID)
	}
	var repos []*model.Repository
	query, args, err := sqlx.In(`SELECT * FROM repositories WHERE user_id IN (?);`, userIDs)
	if err != nil {
		return users, err
	}
	if err := db.Select(&repos, db.Rebind(query), args...); err != nil {
		return users, err
	}
	for _, u := range users {
		for _, r := range repos {
			if r.UserID == *u.ID {
				u.Repos = append(u.Repos, r)
			}
		}
	}

	return users, nil
}

// TODO fix query
func (db *Database) ListRepos() ([]*model.Repository, error) {
	var repos []*model.Repository
	if err := db.Select(&repos, `SELECT * FROM repositories ORDER BY created_at ASC;`); err != nil {
		return repos, err
	}
	userIDs := []int64{}
	for _, r := range repos {
		userIDs = append(userIDs, r.UserID)
	}
	var users []*model.User
	query, args, err := sqlx.In(`SELECT * FROM users WHERE id IN (?);`, userIDs)
	if err != nil {
		return repos, err
	}
	if err := db.Select(&users, db.Rebind(query), args...); err != nil {
		return repos, err
	}
	for _, r := range repos {
		for _, u := range users {
			if r.UserID == *u.ID {
				r.User = u
			}
		}
	}
	return repos, nil
}
