package sqlite

import (
	"context"
	"database/sql"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

type UserRepo struct {
	storage *Store
}

func (r *UserRepo) Add(ctx context.Context, u *store.User) error {
	_, err := r.storage.db.ExecContext(ctx, `INSERT INTO users (id, name) VALUES (?, ?);`,
		u.ChatID,
		u.UserName,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) ByID(ctx context.Context, ID int64) (*store.User, error) {
	u := &store.User{}
	err := r.storage.db.QueryRowContext(ctx,
		`SELECT id, name 
			FROM users
		WHERE id=?;`,
		ID,
	).Scan(&u.ChatID, &u.UserName)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) ByUserName(ctx context.Context, userName string) (*store.User, error) {
	u := &store.User{}
	err := r.storage.db.QueryRowContext(ctx,
		`SELECT id, name FROM users WHERE name=?;`,
		userName,
	).Scan(
		&u.ChatID,
		&userName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NoExistUserError
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) IsExist(ctx context.Context, u *store.User) (bool, error) {
	id := -1
	err := r.storage.db.QueryRowContext(ctx, `SELECT id FROM users WHERE id=?;`, u.ChatID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
