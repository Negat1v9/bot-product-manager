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

// get list products with id
// func (r *ProductListRepo) Get(ctx context.Context, listID int) (*store.ProductList, error) {
// 	stmt, err := r.storage.db.Prepare(
// 		`SELECT (id, owner_id, group_id, name) FROM product_list WHERE id = ?;`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var productList store.ProductList
// 	err = stmt.QueryRowContext(ctx, listID).Scan(
// 		&productList.ID,
// 		&productList.OwnerID,
// 		&productList.GroupID,
// 		&productList.Name,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, store.NoRowListOfProduct
// 		}
// 		return nil, err
// 	}
// 	return &productList, nil
// }

func (r *ProductListRepo) GetListID(ctx context.Context, listName string) (int, error) {
	var id int
	err := r.storage.db.QueryRowContext(
		ctx,
		`SELECT id FROM product_list WHERE name=?`,
		listName,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ProductListRepo) GetAll(ctx context.Context, UserID int) ([]store.ProductList, error) {
	lists := []store.ProductList{}
	row, err := r.storage.db.QueryContext(ctx, `SELECT * FROM product_list WHERE owner_id=?`, UserID)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		var list store.ProductList
		err := row.Scan(&list.ID, &list.OwnerID, &list.GroupID, &list.Name)
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

// delete full list with all product inside
func (r *ProductListRepo) Delete(ctx context.Context, listID int) error {
	stmt, err := r.storage.db.Prepare(
		`DELETE FROM TABLE product_list WHERE id = ?;`)
	if err != nil {
		return err
	}
	if _, err = stmt.QueryContext(ctx, listID); err != nil {
		return err
	}
	return nil
}
