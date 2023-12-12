package sqlite

import (
	"context"
	"database/sql"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

type ProductRepo struct {
	storage *Store
}

func (r *ProductRepo) Add(prod []store.Product) error {
	stmt, err := r.storage.db.Prepare(`INSERT INTO product (name, user_id, list_id) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, p := range prod {
		_, err = stmt.Exec(p.Product, p.UserID, p.ListID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ProductRepo) CountByListID(c context.Context, listID int) (int, error) {
	var count int
	err := r.storage.db.QueryRowContext(c,
		`SELECT count(id) 
			FROM product 
		WHERE product.list_id=?;`,
		listID,
	).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, store.NoProuductExistError
		}
		return 0, err
	}
	return count, nil
}

func (r *ProductRepo) GetByListID(c context.Context, listID, offset, limit int) ([]store.Product, error) {
	rows, err := r.storage.db.QueryContext(c, `
		SELECT id, name, user_id 
			FROM product 
		WHERE product.list_id=? 
		ORDER BY product.id
			LIMIT ? OFFSET ?;`,
		listID,
		limit,
		offset,
	)

	if err != nil {
		return []store.Product{}, err
	}

	defer rows.Close()

	var products = []store.Product{}
	for rows.Next() {
		var prod store.Product
		err = rows.Scan(&prod.ID, &prod.Product, &prod.UserID)
		if err != nil {
			return []store.Product{}, err
		}
		products = append(products, prod)
	}
	return products, nil
}

// func (r *ProductRepo) GetProdWithUserInfo(context.Context)

func (r *ProductRepo) Delete(c context.Context, id int) error {
	_, err := r.storage.db.ExecContext(c, `
	DELETE FROM product WHERE id=?;`, id)
	if err != nil {
		return err
	}
	return err
}
