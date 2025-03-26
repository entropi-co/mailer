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

func (s *Storage) QueryUserByLocal(local string) (*User, error) {
	rows, err := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("id", "local", "created_at").
		From("users").
		Where(squirrel.Eq{"local": local}).
		RunWith(s.Database).
		Query()
	if err != nil {
		return nil, err
	}

	rows.Next()
	var user User
	if err := rows.Scan(&user.ID, &user.Local, &user.CreatedAt); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) QueryUsersByLocals(locals []string) ([]*User, error) {
	rows, err := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("id", "local", "created_at").
		From("users").
		Where(squirrel.Expr("local in ($1)", locals)).
		Query()
	if err != nil {
		return nil, err
	}

	users := make([]*User, len(locals))
	i := 0
	for rows.Next() {
		if err := rows.Scan(&users[i].ID, &users[i].Local, &users[i].CreatedAt); err != nil {
			return nil, err
		}

		i++
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Storage) QueryUserIDsByLocals(locals []string) ([]uint64, error) {
	rows, err := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("id").
		From("users").
		Where(squirrel.Expr("local in ($1)", locals)).
		Query()
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, len(locals))
	i := 0
	for rows.Next() {
		if err := rows.Scan(&ids[i]); err != nil {
			return nil, err
		}

		i++
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *Storage) QueryUserByLocalAndKeyValue(local string, value string) (*User, error) {
	builder := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar)

	row := builder.
		Select("id", "local", "created_at").
		From("users").
		LeftJoin("keys k ON k.owner = users.id").
		Where(squirrel.Eq{"users.local": local, "k.value": value}).
		RunWith(s.Database).
		QueryRow()

	var user User
	if err := row.Scan(&user.ID, &user.Local, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}
