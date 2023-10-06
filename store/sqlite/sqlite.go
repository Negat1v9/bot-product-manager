package sqlite

import (
	"database/sql"
	"errors"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

type Store struct {
	db               *sql.DB
	userRepo         *UserRepo
	managerGroupRepo *ManagerGroupRepo
	groupRepo        *GroupRepo
	productListRepo  *ProductListRepo
	produtRepo       *ProductRepo
}

func Newstorage(db *sql.DB) store.Store {
	return &Store{
		db: db,
	}
}

// TODO: Make for each table function
func (s *Store) CreateTables() error {
	const pkg = "store/sqlite"
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS "users" (
		id INTEGER PRIMARY KEY,
		name TEXT
	);`)
	if err != nil {
		return errors.New(pkg + ": " + err.Error())
	}

	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS "manager_group"(
		id INTEGER PRIMARY KEY,
		owner_id INTEGER NOT NULL,
		FOREIGN KEY (owner_id) REFERENCES users (id)
		);`)
	if err != nil {
		return errors.New(pkg + ": " + err.Error())
	}

	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS "group" (
		user_id INTEGER NOT NULL,
		group_id INTEGER NOT NULL,
		PRIMARY KEY(user_id, group_id),
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		FOREIGN KEY (group_id) REFERENCES manager_group(id) ON DELETE CASCADE
		);`)
	if err != nil {
		return errors.New(pkg + ": " + err.Error())
	}

	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS "product_list"(
		id INTEGER PRIMARY KEY,
		owner_id INTEGER NOT NULL,
		group_id INTEGER DEFAULT NULL,
		name TEXT NOT NULL,
		FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (group_id) REFERENCES manager_group(id) ON DELETE CASCADE

	);`)
	if err != nil {
		return errors.New(pkg + ": " + err.Error())
	}

	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS "product"(
		id INTEGER PRIMARY KEY,
		list_id INTEGER NOT NULL,
		products TEXT DEFALUT ' ',
		FOREIGN KEY (list_id) REFERENCES product_list(id) ON DELETE CASCADE
	);`)
	if err != nil {
		return errors.New(pkg + ": " + err.Error())
	}
	return nil
}

func (s *Store) User() store.UserRepo {
	if s.userRepo != nil {
		return s.userRepo
	}
	s.userRepo = &UserRepo{
		storage: s,
	}
	return s.userRepo
}

func (s *Store) ManagerGroup() store.ManagerGroupRepo {
	if s.managerGroupRepo != nil {
		return s.managerGroupRepo
	}
	s.managerGroupRepo = &ManagerGroupRepo{
		storage: s,
	}
	return s.managerGroupRepo
}

func (s *Store) Group() store.GroupRepo {
	if s.groupRepo != nil {
		return s.groupRepo
	}
	s.groupRepo = &GroupRepo{
		storage: s,
	}
	return s.groupRepo
}

func (s *Store) ProductList() store.ProductListRepo {
	if s.productListRepo != nil {
		return s.productListRepo
	}
	s.productListRepo = &ProductListRepo{
		storage: s,
	}
	return s.productListRepo
}

func (s *Store) Product() store.ProductRepo {
	if s.produtRepo != nil {
		return s.produtRepo
	}
	s.produtRepo = &ProductRepo{
		storage: s,
	}
	return s.produtRepo
}
