package store

import "context"

type Store interface {
	CreateTables() error
	User() UserRepo
	ManagerGroup() ManagerGroupRepo
	Group() GroupRepo
	ProductList() ProductListRepo
	Product() ProductRepo
}

type UserRepo interface {
	Add(ctx context.Context, u *User) error
	IsExist(ctx context.Context, u *User) (bool, error)
}

type ManagerGroupRepo interface {
	Create(ctx context.Context, ownerID int) error
}

type GroupRepo interface {
	AddUser(ctx context.Context, g *Group) error
	DeleteUser(ctx context.Context, g *Group) error
}

// NOTE: Merge ProductList and Product Repositories
type ProductListRepo interface {
	Create(ctx context.Context, p *ProductList) error
	Get(ctx context.Context, litsID int) (*ProductList, error)
	GetAll(ctx context.Context, UserID int) ([]ProductList, error)
	Delete(ctx context.Context, listID int) error
}

type ProductRepo interface {
	Add(ctx context.Context, p *Product) error
	Delete(ctx context.Context, productID int) error
}

type User struct {
	ChatID   int
	UserName string
}

type ManagerGroup struct {
	ID       int
	Owner_id int
}

type Group struct {
	UserID  int
	GroupID int
}

// One list products
type ProductList struct {
	ID      int
	OwnerID int64
	GroupID *int
	Name    string
}

// product copy
type Product struct {
	ID     int
	Name   string
	ListID int
}
