package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

type ProductListRepo struct {
	storage *Store
}

// create list for add there product
func (r *ProductListRepo) Create(ctx context.Context, p *store.ProductList) (id int, err error) {
	var prodStr string
	if len(p.Products) != 0 {
		prodStr, err = r.convertProductString(p.Products)
		if err != nil {
			return 0, err
		}
	}
	var stmt string
	var row *sql.Row
	if p.GroupID != nil {
		stmt = `INSERT INTO product_list (owner_id, group_id, name, products) VALUES (?, ?, ?, ?)
			RETURNING product_list.id;`
		row = r.storage.db.QueryRowContext(ctx, stmt, p.OwnerID, p.GroupID, p.Name, prodStr)
	} else {
		stmt = `INSERT INTO product_list (owner_id, name, products) VALUES (?, ?, ?)
			RETURNING product_list.id;`
		row = r.storage.db.QueryRowContext(ctx, stmt, p.OwnerID, p.Name, prodStr)
	}

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Info: get ID of solo created Product list
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

// get group list
func (r *ProductListRepo) GetGroupListByID(ctx context.Context, listID int) (store.ProductList, error) {
	var p store.ProductList
	err := r.storage.db.QueryRowContext(ctx,
		`SELECT product_list.owner_id,
				product_list.group_id, product_list.name
		 FROM product_list 
		 WHERE product_list.id=? AND is_active=1;`, listID,
	).Scan(&p.OwnerID, &p.GroupID, &p.Name)
	if err != nil {
		return store.ProductList{}, err
	}
	return p, nil
}

// Info: If want receive by id name must by "" empty string
func (r *ProductListRepo) GetAllInfoProductLissIdOrName(ctx context.Context, listID int, name string) (*store.ProductList, error) {
	var prodStr, editStr *string
	var p store.ProductList
	var row *sql.Row

	stmt := `SELECT id, owner_id, group_id, name, products, editors
			FROM product_list `

	if name == "" {
		stmt += "WHERE product_list.id=?;"
		row = r.storage.db.QueryRowContext(ctx, stmt, listID)
	} else {
		stmt += "WHERE product_list.name=?;"
		row = r.storage.db.QueryRowContext(ctx, stmt, name)
	}

	err := row.Scan(&p.ID, &p.OwnerID, &p.GroupID, &p.Name, &prodStr, &editStr)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NoRowProductError
		}
		return nil, err
	}
	if prodStr != nil {
		p.Products, err = r.convertStringProduct(*prodStr)
		if err != nil {
			return nil, err
		}
	}
	if editStr != nil {
		p.Editors, err = r.convertStringEditors(*editStr)
		if err != nil {
			return nil, err
		}
	}
	return &p, nil
}

func (r *ProductListRepo) EditProductList(ctx context.Context, l store.ProductList) error {
	productsString, err := r.convertProductString(l.Products)

	if err != nil {
		return err
	}

	editorsString, err := r.convertEditorsString(l.Editors)
	if err != nil {
		return err
	}
	stmt := `UPDATE product_list 
				SET products=?, editors=? 
			WHERE id=? 
			RETURNING id;`
	var id *int
	err = r.storage.db.QueryRowContext(ctx,
		stmt,
		productsString,
		editorsString,
		l.ID,
	).Scan(&id)
	if err != nil {
		return err
	}
	if id != nil {
		return nil
	}
	return store.NoRowExistError
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

func (r *ProductListRepo) convertStringProduct(s string) ([]string, error) {
	items := make([]string, 0)
	err := json.Unmarshal([]byte(s), &items)
	if err != nil {
		return []string{}, nil
	}
	return items, nil
}
func (r *ProductListRepo) convertProductString(p []string) (string, error) {
	res, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (r *ProductListRepo) convertStringEditors(s string) ([]store.Editors, error) {
	editors := make([]store.Editors, 0)
	err := json.Unmarshal([]byte(s), &editors)
	if err != nil {
		return []store.Editors{}, nil
	}
	return editors, nil
}

func (r *ProductListRepo) convertEditorsString(e []store.Editors) (string, error) {
	res, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
