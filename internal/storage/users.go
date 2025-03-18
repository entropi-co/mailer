package storage

import (
	"github.com/Masterminds/squirrel"
	"time"
)

type User struct {
	ID        uint64    `json:"id"`
	Local     string    `json:"local"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Storage) CreateUser(id uint64) (*User, error) {
	rows, err := squirrel.
		Insert("users").
		Columns("id").
		Values(id).
		Suffix("RETURNING id, created_at").
		RunWith(s.Database).
		Query()
	if err != nil {
		return nil, err
	}

	var user User

	rows.Next()
	if err := rows.Scan(&user.ID, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) QueryUsersByLocals(locals []string) ([]*User, error) {
	return nil, nil
}

func (s *Storage) QueryUserIDsByLocals(id uint64) ([]uint64, error) {
	return nil, nil
}
