package sqlite

import "context"

// Repo for Group
type ManagerGroupRepo struct {
	storage *Store
}

// Create group
func (r *ManagerGroupRepo) Create(ctx context.Context, ownerID int) error {
	_, err := r.storage.db.ExecContext(
		ctx,
		`INSERT INTO manager_group (owner_id) VALUES (?);`,
		ownerID)
	if err != nil {
		return err
	}
	return nil
}
