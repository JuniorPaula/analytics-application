package database

import (
	"c2d-reports/internal/config"
	"database/sql"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.ConnectStringDatabase)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
