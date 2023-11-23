package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

type ProductRepo struct {
	storage *Store
}

func (r *ProductRepo) Create(ctx context.Context, p *store.Product) (err error) {
	var prodStr string
	if len(p.Products) != 0 {
		prodStr, err = r.ConvertProductString(p.Products)
		if err != nil {
			return err
		}
	}
	_, err = r.storage.db.ExecContext(ctx,
		`INSERT INTO product (list_id, products) VALUES (?, ?);`,
		p.ListID, prodStr,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) Add(ctx context.Context, p store.Product) error {
	productsString, err := r.ConvertProductString(p.Products)

	if err != nil {
		return err
	}

	editorsString, err := r.ConvertEditorsString(p.Editors)
	if err != nil {
		return nil
	}
	stmt := `UPDATE product 
				SET products=?, editors=? 
			WHERE list_id=? 
			RETURNING id;`
	var id *int
	err = r.storage.db.QueryRowContext(ctx,
		stmt,
		productsString,
		editorsString,
		p.ListID,
	).Scan(&id)
	if err != nil {
		return err
	}
	if id != nil {
		return nil
	}
	return store.NoRowExistError
}

// Info: withText: Boll - parametr for getting in Obj original string from DB
func (r *ProductRepo) GetAllProducts(ctx context.Context, listID int) (*store.Product, error) {
	var productsStr string
	var editorsStr string
	var p store.Product
	err := r.storage.db.QueryRowContext(ctx,
		`SELECT product.id, product.list_id, product.products, product.editors,
			product_list.name 
		FROM product 
		JOIN product_list
			ON product.list_id = product_list.id
		WHERE product.list_id=?;`,
		listID,
	).Scan(&p.ID, &p.ListID, &productsStr, &editorsStr, &p.ListName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NoRowProductError
		}
		return nil, err
	}

	p.Products, err = r.ConvertStringProduct(productsStr)
	if err != nil {
		return nil, err
	}
	p.Editors, err = r.ConvertStringEditors(editorsStr)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepo) Delete(ctx context.Context, ID int) error {
	_, err := r.storage.db.ExecContext(
		ctx,
		`DELETE FROM product WHERE list_id=?;`,
		ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepo) ConvertStringProduct(s string) ([]string, error) {
	items := make([]string, 0)
	err := json.Unmarshal([]byte(s), &items)
	if err != nil {
		return []string{}, nil
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

func (r *ProductRepo) ConvertStringEditors(s string) ([]store.Editors, error) {
	editors := make([]store.Editors, 0)
	err := json.Unmarshal([]byte(s), &editors)
	if err != nil {
		return []store.Editors{}, nil
	}
	return editors, nil
}

func (r *ProductRepo) ConvertEditorsString(e []store.Editors) (string, error) {
	res, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
