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

func (db *Database) ListUsers() ([]*model.User, error) {
	var users []*model.User
	if err := db.Select(&users, `SELECT * FROM users ORDER BY created_at ASC;`); err != nil {
		return users, err
	}
	return users, nil
}
