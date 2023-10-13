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
	Create(ctx context.Context, g *GroupInfo) (int, error)
	ByGroupName(ctx context.Context, groupName string) (*GroupInfo, error)
	ByGroupID(ctx context.Context, groupID int) (*GroupInfo, error)
	AllByGroupID(ctx context.Context, groupID int) (*GroupList, error)
	UserGroup(ctx context.Context, userID int) ([]GroupInfo, error)
	InfoGroup(ctx context.Context, groupID int) (*GroupInfo, error)
}

type GroupRepo interface {
	AddUser(ctx context.Context, g *Group) error
	DeleteUser(ctx context.Context, g *Group) error
}

// NOTE: Merge ProductList and Product Repositories
type ProductListRepo interface {
	Create(ctx context.Context, p *ProductList) error
	GetListID(ctx context.Context, listName string) (int, error)
	GetAll(ctx context.Context, UserID int) ([]ProductList, error)
	Delete(ctx context.Context, listID int) error
}

type ProductRepo interface {
	Create(ctx context.Context, listID int) error
	GetAll(ctx context.Context, listID int) (*Product, error)
	Add(ctx context.Context, p Product) error
	Delete(ctx context.Context, productID int) error
}

type User struct {
	ChatID   int
	UserName string
}

type GroupInfo struct {
	ID        int
	OwnerID   int64
	GroupName string
	UsersInfo *[]User
}

type Group struct {
	UserID  int64
	GroupID int
}

// One list products all value is a poinet ->
// group can have 0 products list
type ProductList struct {
	ID      *int
	OwnerID *int64
	GroupID *int
	Name    *string
}

// product copy
type Product struct {
	ID       int
	ListID   int
	Products []string
}

type GroupList struct {
	PruductLists []ProductList
	GroupOwnerID int64
}
