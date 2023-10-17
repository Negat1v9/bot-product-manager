package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

const (
	defaultValueProducts = `"[]"`
)

type ProductRepo struct {
	storage *Store
}

func (r *ProductRepo) Create(ctx context.Context, listID int) error {
	_, err := r.storage.db.ExecContext(ctx,
		`INSERT INTO product (list_id, products) VALUES (?, ?);`,
		listID,
		defaultValueProducts,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) Add(ctx context.Context, p store.Product) error {
	productsString, err := r.ConvertProductString(p.Products)
	fmt.Println("product-repo", productsString)
	if err != nil {
		return err
	}
	var id *int
	stmt := `UPDATE product
			SET products=? 
			WHERE list_id=?
			RETURNING id;`

	err = r.storage.db.QueryRowContext(ctx, stmt, productsString, p.ListID).Scan(&id)
	if err != nil {
		return err
	}
	if id == nil {
		fmt.Println("products update was not")
	}
	return nil
}

func (r *ProductRepo) GetAll(ctx context.Context, listID int) (*store.Product, error) {
	var products string
	var p store.Product
	err := r.storage.db.QueryRowContext(ctx,
		`SELECT id, list_id, products FROM product WHERE list_id=?`,
		listID,
	).Scan(&p.ID, &p.ListID, &products)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NoRowProductError
		}
		return nil, err
	}
	p.Products, err = r.ConvertStringProduct(products)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepo) Delete(ctx context.Context, productID int) error {
	_, err := r.storage.db.ExecContext(
		ctx,
		`DELETE FROM product WHERE id=?;`,
		productID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) ConvertStringProduct(s string) ([]string, error) {
	var items []string
	err := json.Unmarshal([]byte(s), &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
func (r *ProductRepo) ConvertProductString(p []string) (string, error) {
	res, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
