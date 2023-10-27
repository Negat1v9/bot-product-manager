package sqlite

import (
	"context"
	"github.com/Negat1v9/telegram-bot-orders/store"
)

type ProductListRepo struct {
	storage *Store
}

// create list for add there product
func (r *ProductListRepo) Create(ctx context.Context, p *store.ProductList) error {
	stmt, err := r.storage.db.Prepare(`
	INSERT INTO product_list (owner_id, group_id, name) VALUES (?, ?, ?);`)
	if err != nil {
		return err
	}
	if _, err = stmt.ExecContext(ctx, p.OwnerID, p.GroupID, p.Name); err != nil {
		return err
	}

	return nil
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
		 FROM product_list WHERE owner_id=? AND group_id IS NULL;`, userID)
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
