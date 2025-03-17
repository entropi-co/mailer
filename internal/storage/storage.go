package storage

import (
	"database/sql"
)

type Storage struct {
	Database *sql.DB
}

func CreateStorage() *Storage {
	return &Storage{Database: ConnectDatabase()}
}
