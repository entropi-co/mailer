package storage

import "github.com/Masterminds/squirrel"

type Key struct {
	/// Value is actual api key value
	Value string `json:"value"`
	/// Owner refers to users
	Owner uint64 `json:"owner"`
}

func (s *Storage) QueryKeyByLocal(local string) (*Key, error) {
	// select * from keys where owner = (select id from users where local = $1);

	builder := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar)

	row := builder.
		Select("value", "owner").
		From("keys").
		Where(
			builder.
				Select("id").
				From("users").
				Where(squirrel.Eq{"local": local}).
				Prefix("owner IN (").
				Suffix(")"),
		).
		RunWith(s.Database).
		QueryRow()

	var key Key
	if err := row.Scan(&key.Value, &key.Owner); err != nil {
		return nil, err
	}

	return &key, nil
}
