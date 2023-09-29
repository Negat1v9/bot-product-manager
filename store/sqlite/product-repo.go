package sqlite

import (
	"context"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

type ProductRepo struct {
	storage *Store
}

func (r *ProductRepo) Add(ctx context.Context, p *store.Product) error {
	_, err := r.storage.db.ExecContext(
		ctx,
		`INSERT INTO product (name, list_id) VALUES (?, ?);`,
		p.Name,
		p.ListID,
	)
	if err != nil {
		return err
	}

	return nil
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
