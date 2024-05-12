package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type Storage struct {
	db *sql.DB
}

var embedMigrations embed.FS

func NewConnect(storagePath string) (*Storage, error) {
	const op = "Storage.sqlite.NewConnect"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// проверка подключения
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = goose.Up(db, "./db/migrations"); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(SaveURL string, alias string) (*int64, error) {
	const op = "Storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO urls(url, alias) values(?,?)")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(SaveURL, alias)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &id, nil
}

func (s *Storage) GetAlias(LoadURL string) (*string, error) {
	const op = "Storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT alias FROM urls WHERE url=?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var alias string

	err = stmt.QueryRow(LoadURL).Scan(&alias)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &alias, nil
}

func (s *Storage) DeleteURL(LoadURL string) error {
	const op = "Storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare(`DELETE FROM urls WHERE url=?`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(LoadURL)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
