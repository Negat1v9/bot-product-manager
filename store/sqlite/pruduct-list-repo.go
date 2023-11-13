package sqlite

import (
	"context"
	"github.com/Negat1v9/telegram-bot-orders/store"
)

type ProductListRepo struct {
	storage *Store
}

// create list for add there product
func (r *ProductListRepo) Create(ctx context.Context, p *store.ProductList) (int, error) {
	stmt, err := r.storage.db.Prepare(`
	INSERT INTO product_list (owner_id, group_id, name) VALUES (?, ?, ?)
		RETURNING product_list.id;`)
	if err != nil {
		return 0, err
	}
	var id int
	err = stmt.QueryRowContext(ctx, p.OwnerID, p.GroupID, p.Name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ProductListRepo) GetListID(ctx context.Context, listName string) (int, error) {
	var id int
	err := r.storage.db.QueryRowContext(
		ctx,
		`SELECT id FROM product_list WHERE name=?;`,
		listName,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ProductListRepo) GetAll(ctx context.Context, userID int64) ([]store.ProductList, error) {
	lists := []store.ProductList{}
	row, err := r.storage.db.QueryContext(ctx,
		`SELECT product_list.id, product_list.owner_id, product_list.group_id, product_list.name
		 FROM product_list WHERE owner_id=? AND group_id IS NULL AND is_template IS NULL;`, userID)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		var list store.ProductList
		err = row.Scan(&list.ID, &list.OwnerID, &list.GroupID, &list.Name)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	if len(lists) == 0 {
		return nil, store.NoRowListOfProductError
	}
	return lists, nil
}

func (r *ProductListRepo) GetAllGroupTemplates(ctx context.Context, groupID int) ([]store.ProductList, error) {
	rows, err := r.storage.db.QueryContext(ctx,
		`SELECT product_list.id, product_list.owner_id, product_list.group_id, product_list.name
			FROM product_list
			JOIN group_info
			ON product_list.group_id = group_info.id
		WHERE group_info.id=? AND product_list.is_template IS NOT NULL;`,
		groupID,
	)
	if err != nil {
		return nil, err
	}
	lists := []store.ProductList{}
	for rows.Next() {
		list := store.ProductList{}
		err = rows.Scan(&list.ID, &list.OwnerID, &list.GroupID, &list.Name)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	return lists, nil
}

func (r *ProductListRepo) MergeListGroup(ctx context.Context, listID, groupID int) error {
	_, err := r.storage.db.ExecContext(ctx,
		`UPDATE product_list
			SET group_id=?
		WHERE product_list.id=?;`,
		groupID,
		listID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductListRepo) SaveListAsTemplate(ctx context.Context, listID int) error {
	_, err := r.storage.db.ExecContext(ctx,
		`UPDATE product_list
			SET is_template="TRUE"
		WHERE product_list.id=?;`,
		listID,
	)
	if err != nil {
		return err
	}
	return nil
}

// delete full list with all product inside
func (r *ProductListRepo) Delete(ctx context.Context, listID int) error {
	_, err := r.storage.db.ExecContext(ctx,
		`DELETE FROM product_list WHERE id = ?;`,
		listID,
	)
	if err != nil {
		return err
	}
	return nil
}
