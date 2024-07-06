package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func New(dbPath string) (*Database, error) {
	sqliteDatabase, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return &Database{}, err
	}
	db := &Database{
		db: sqliteDatabase,
	}
	if err := db.CreateTables(); err != nil {
		return db, err
	}
	return db, nil
}

func (s *Database) CreateTables() error {
	createTablesSQL := `
CREATE TABLE IF NOT EXISTS users(id integer primary key autoincrement not null, name text, pass text, ssh text);
INSERT INTO users VALUES(1, 'root', 'password', 'some-ssh-key', 'keykeys');
CREATE TABLE repositories(id integer primary key autoincrement not null, path text, foreign key(id,name) references repositories ON DELETE CASCADE ON UPDATE CASCADE);`

	statement, err := s.db.Prepare(createTablesSQL)
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
