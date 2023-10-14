package sqlite

import (
	"context"
	"database/sql"

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

// INFO: this functinal return info of group and usersInfo if he is in group
// make sure about user in group or not, and get all parametrs to make all you need
// with it
func (r *GroupRepo) GroupByUserAndGroupID(ctx context.Context, userID int64, groupID int) (*store.GroupInfo, error) {
	groupInfo := &store.GroupInfo{}
	user := &store.User{}
	err := r.storage.db.QueryRowContext(ctx,
		`SELECT user_group.user_id,
			users.name,
			group_info.id, group_info.group_name, group_info.owner_id
		FROM user_group
		JOIN users
		ON users.id = user_group.user_id
			JOIN group_info
			ON group_info.id = user_group.group_id
		WHERE user_group.user_id=? AND user_group.group_id=?;`,
		userID, groupID,
	).Scan(&user.ChatID, &user.UserName,
		&groupInfo.ID, &groupInfo.GroupName, &groupInfo.OwnerID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	groupInfo.UsersInfo = &[]store.User{*user}
	return groupInfo, nil
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
