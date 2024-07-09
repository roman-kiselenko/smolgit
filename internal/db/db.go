package db

import (
	_ "embed"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db     *sqlx.DB
	logger *slog.Logger
}

//go:embed schema.sql
var schema string

func New(logger *slog.Logger, dbPath string) (*Database, error) {
	sqliteDatabase, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return &Database{}, err
	}
	db := &Database{
		db:     sqliteDatabase,
		logger: logger,
	}
	if _, err := db.db.MustExec(schema).RowsAffected(); err != nil {
		logger.Error("migrate database scheme", "result", err)
	}
	// var root int
	// _ = db.db.Get(&root, `SELECT 1 FROM users WHERE email = $1;`, "root")
	// if root == 0 {
	// 	hpass, _ := password.HashPassword("password")
	// 	if _, err = db.db.NamedExec(`INSERT INTO users (email, password) VALUES (:email, :password)`,
	// 		map[string]interface{}{
	// 			"email":    "root",
	// 			"password": hpass,
	// 		}); err != nil {
	// 		logger.Error("cant insert user", "err", err)
	// 	}
	// 	logger.Debug("root user created", "user", "root")

	// }
	return db, nil
}

func (db *Database) FindUserFromKey(key, user string) (string, error) {
	db.logger.Debug("fetch user by key", "key", key[:10], "user", user)
	return "", nil
}
