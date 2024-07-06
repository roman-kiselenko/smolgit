package db

import (
	"git-pebble/internal/password"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sqlx.DB
}

var schema = `
CREATE TABLE IF NOT EXISTS users(
    id         INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	email      VARCHAR(250) DEFAULT '',
	password   VARCHAR(250) DEFAULT NULL
);
CREATE TABLE IF NOT EXISTS repositories(
    id         INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	path       VARCHAR(250) DEFAULT '',
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);`

func New(dbPath string) (*Database, error) {
	sqliteDatabase, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return &Database{}, err
	}
	db := &Database{
		db: sqliteDatabase,
	}
	db.db.MustExec(schema)
	var root int
	_ = db.db.Get(&root, `SELECT 1 FROM users WHERE email = $1;`, "root")
	if root == 0 {
		hpass, _ := password.HashPassword("password")
		_, err = db.db.NamedExec(`INSERT INTO users (email, password) VALUES (:email, :password)`,
			map[string]interface{}{
				"email":    "root",
				"password": hpass,
			})
	}
	return db, err
}

func (s *Database) CreateTables(table string) error {
	statement, err := s.db.Prepare(table)
	if err != nil {
		return err
	}
	if _, err := statement.Exec(); err != nil {
		return err
	}
	return nil
}

func (s *Database) Close() error {
	return s.db.Close()
}
