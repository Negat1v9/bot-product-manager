package sqlite

import (
	"context"
	"database/sql"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

// Repo for Group
type ManagerGroupRepo struct {
	storage *Store
}

// Create group
func (r *ManagerGroupRepo) Create(ctx context.Context, g *store.GroupInfo) (int, error) {
	var id int
	err := r.storage.db.QueryRowContext(
		ctx,
		`INSERT INTO group_info (owner_id, group_name) VALUES (?, ?) RETURNING id;`,
		g.OwnerID,
		g.GroupName,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ManagerGroupRepo) ByGroupID(ctx context.Context, groupID int) (*store.GroupInfo, error) {
	groupInfo := store.GroupInfo{}
	err := r.storage.db.QueryRowContext(ctx,
		`SELECT group_info.id, group_info.owner_id, group_info.group_name
		FROM group_info
		WHERE id=?;`,
		groupID).Scan(
		&groupInfo.ID,
		&groupInfo.OwnerID,
		&groupInfo.GroupName,
	)
	if err != nil {
		return nil, err
	}
	return &groupInfo, nil
}

func (r *ManagerGroupRepo) ByGroupName(ctx context.Context, gN string) (*store.GroupInfo, error) {
	groupInfo := store.GroupInfo{}
	err := r.storage.db.QueryRowContext(ctx,
		`SELECT group_info.id, group_info.group_name, group_info.owner_id
			FROM group_info
		WHERE group_info.group_name = ?;`,
		gN,
	).Scan(&groupInfo.ID,
		&groupInfo.GroupName,
		&groupInfo.OwnerID,
	)

	if err != nil {
		return nil, err
	}

	return &groupInfo, nil
}

// FIXME: move to product-list repository
func (r *ManagerGroupRepo) AllByGroupID(ctx context.Context, groupID int) (*store.GroupList, error) {
	groupList := store.GroupList{}
	rows, err := r.storage.db.QueryContext(ctx,
		`SELECT product_list.id, product_list.owner_id, product_list.name, group_info.owner_id
			FROM group_info
			LEFT JOIN product_list 
			ON group_info.id = product_list.group_id
		WHERE group_info.id = ? AND product_list.is_template IS NULL;`,
		groupID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &groupList, store.NoRowListOfProductError
		}
		return nil, err
	}
	defer rows.Close()
	var groupOwnerID int64
	for rows.Next() {
		var list store.ProductList
		err = rows.Scan(&list.ID, &list.OwnerID, &list.Name, &groupOwnerID)
		if err != nil {
			return nil, err
		}
		// take one value from list to be confirm its not nil
		if list.ID != nil {
			groupList.PruductLists = append(groupList.PruductLists, list)
		} else {
			break
		}
	}
	groupList.GroupOwnerID = groupOwnerID
	if len(groupList.PruductLists) == 0 {
		return &groupList, store.NoRowListOfProductError
	}
	return &groupList, nil
}

// FIXME: move to group-repository
func (r *ManagerGroupRepo) InfoGroup(ctx context.Context, groupID int) (*store.GroupInfo, error) {
	rows, err := r.storage.db.QueryContext(ctx,
		`SELECT group_info.id, group_info.group_name, group_info.owner_id,
			users.name,
			user_group.user_id
		FROM user_group
		JOIN group_info
		ON user_group.group_id = group_info.id
			JOIN users 
			ON user_group.user_id = users.id
	WHERE user_group.group_id=?;`,
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var grInfo store.GroupInfo
	var users []store.User
	for rows.Next() {
		user := store.User{}
		if err = rows.Scan(&grInfo.ID, &grInfo.GroupName, &grInfo.OwnerID,
			&user.UserName, &user.ChatID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	grInfo.UsersInfo = &users
	return &grInfo, nil
}

// InfoGroupByName
// FIXME: move to group-repository
func (r *ManagerGroupRepo) UserGroup(ctx context.Context, userID int64) ([]store.GroupInfo, error) {
	rows, err := r.storage.db.QueryContext(ctx,
		`SELECT group_info.id, group_info.group_name, group_info.owner_id
		FROM user_group 
		LEFT JOIN group_info
		ON user_group.group_id = group_info.id 
	WHERE user_group.user_id = ?;`,
		userID,
	)
	var groups = []store.GroupInfo{}
	if err != nil {
		if err == sql.ErrNoRows {
			return groups, store.NoUserGroupError
		}
		return nil, err
	}
	var group store.GroupInfo
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&group.ID, &group.GroupName, &group.OwnerID)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	if len(groups) == 0 {
		return groups, store.NoUserGroupError
	}

	return groups, nil
}

func (r *ManagerGroupRepo) DeleteGroup(ctx context.Context, groupID int) error {
	_, err := r.storage.db.ExecContext(ctx,
		`DELETE FROM group_info
		WHERE group_info.id=?;`,
		groupID,
	)
	if err != nil {
		return err
	}
	return nil
}
