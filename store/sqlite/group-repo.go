package sqlite

import (
	"context"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

type GroupRepo struct {
	storage *Store
}

func (r *GroupRepo) AddUser(ctx context.Context, g *store.Group) error {
	_, err := r.storage.db.ExecContext(
		ctx,
		`INSERT INTO user_group (user_id, group_id) VALUES (?, ?);`,
		g.UserID,
		g.GroupID,
	)
	if err != nil {
		return err
	}

	return nil
}

// For delete user from group needs to user exist in this group
func (r *GroupRepo) DeleteUser(ctx context.Context, g *store.Group) error {
	_, err := r.storage.db.ExecContext(
		ctx,
		`DELETE FROM user_group WHERE user_id=? AND group_id=?;`,
		g.UserID,
		g.GroupID,
	)
	if err != nil {
		return err
	}
	return nil
}
