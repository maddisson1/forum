package database

import (
	"database/sql"
	"os"
	"strings"
)

const (
	migrate = "./migrate/create.sql"
)

type Storage struct {
	DB *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}

	if err = Create(db); err != nil {
		return nil, err
	}
	// Создаем индекс, если он не существует
	// _, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_post_created ON posts(created)`)
	// if err != nil {
	// 	return nil, err
	// }

	return &Storage{DB: db}, nil
}

func Create(db *sql.DB) error {
	file, err := os.ReadFile(migrate)
	if err != nil {
		return err
	}
	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		_, err := db.Exec(request)
		if err != nil {
			return err
		}
	}
	return nil
}
