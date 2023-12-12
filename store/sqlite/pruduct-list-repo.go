package sqlite

import (
	"context"
	"database/sql"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

type ProductListRepo struct {
	storage *Store
}

// create list for add there product

func (r *ProductListRepo) Create(ctx context.Context, p *store.ProductList) (id int, err error) {

	var stmt string
	var row *sql.Row
	if p.GroupID != nil {
		stmt = `INSERT INTO product_list (owner_id, group_id, name) VALUES (?, ?, ?)
			RETURNING product_list.id;`
		row = r.storage.db.QueryRowContext(ctx, stmt, p.OwnerID, p.GroupID, p.Name)
	} else {
		stmt = `INSERT INTO product_list (owner_id, name) VALUES (?, ?)
			RETURNING product_list.id;`
		row = r.storage.db.QueryRowContext(ctx, stmt, p.OwnerID, p.Name)
	}

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Info: get solo create listlists
func (r *ProductListRepo) GetAllNames(ctx context.Context, userID int64) ([]store.ProductList, error) {
	lists := []store.ProductList{}
	// var sProd, sEditors *string
	row, err := r.storage.db.QueryContext(ctx,
		`SELECT product_list.id, product_list.owner_id, product_list.name
		 FROM product_list 
		 WHERE owner_id=? AND group_id IS NULL AND is_active=1;`, userID)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		var list store.ProductList
		err = row.Scan(&list.ID, &list.OwnerID, &list.Name)
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

// Info: If want receive by id name must by "" empty string
func (r *ProductListRepo) GetFoolInfoGroupProdList(ctx context.Context, listID int) (*store.FoolInfoProductList, error) {
	rows, err := r.storage.db.QueryContext(ctx,
		`SELECT product.name, 
			product_list.id, product_list.group_id, product_list.name,
			users.name
			FROM product
		JOIN product_list ON product.list_id = product_list.id
		JOIN users ON product.user_id = users.id
			WHERE product.list_id=?
		ORDER BY product.id;`,
		listID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NoRowProductError
		}
		return nil, err
	}
	defer rows.Close()
	var p store.FoolInfoProductList
	for rows.Next() {

		prod := store.InfoProduct{}
		user := &store.User{}
		err = rows.Scan(&prod.Product,
			&p.List.ID, &p.List.GroupID, &p.List.Name,
			&user.UserName,
		)
		if err != nil {
			return nil, err
		}
		prod.User = user
		p.Products = append(p.Products, prod)
	}
	return &p, nil
}

func (r *ProductListRepo) GetFoolInfoProdList(ctx context.Context, listID int) (*store.FoolInfoProductList, error) {
	rows, err := r.storage.db.QueryContext(ctx,
		`SELECT product.name, 
			product_list.id, product_list.name
			FROM product
		LEFT JOIN product_list ON product.list_id = product_list.id
			WHERE product.list_id=?
		ORDER BY product.id;`,
		listID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &store.FoolInfoProductList{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	var p store.FoolInfoProductList
	for rows.Next() {

		prod := store.InfoProduct{}
		err = rows.Scan(&prod.Product,
			&p.List.ID, &p.List.Name,
		)
		if err != nil {
			return nil, err
		}
		p.Products = append(p.Products, prod)
	}
	return &p, nil
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

func (r *ProductListRepo) MakeListInactive(ctx context.Context, listID int) error {
	_, err := r.storage.db.ExecContext(ctx,
		`UPDATE product_list
			SET is_active=0
		WHERE product_list.id=?;`,
		listID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductListRepo) MakeListActive(ctx context.Context, listID int) (string, error) {
	var listName string
	err := r.storage.db.QueryRowContext(ctx,
		`UPDATE product_list
			SET is_active=1
		WHERE product_list.id=?
		RETURNING product_list.name;`,
		listID,
	).Scan(&listName)
	if err != nil {
		return "", err
	}
	return listName, nil
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
