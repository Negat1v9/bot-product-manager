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
	ByID(ctx context.Context, ID int64) (*User, error)
	ByUserName(ctx context.Context, userName string) (*User, error)
	IsExist(ctx context.Context, u *User) (bool, error)
}

type ManagerGroupRepo interface {
	Create(ctx context.Context, g *GroupInfo) (int, error)
	ByGroupName(ctx context.Context, groupName string) (*GroupInfo, error)
	ByGroupID(ctx context.Context, groupID int) (*GroupInfo, error)
	AllByGroupID(ctx context.Context, groupID int) (*GroupList, error)
	UserGroup(ctx context.Context, userID int64) ([]GroupInfo, error)
	InfoGroup(ctx context.Context, groupID int) (*GroupInfo, error)
	DeleteGroup(ctx context.Context, groupID int) error
}

type GroupRepo interface {
	AddUser(ctx context.Context, g *Group) error
	DeleteUser(ctx context.Context, g *Group) error
}

// NOTE: Merge ProductList and Product Repositories
type ProductListRepo interface {
	Create(ctx context.Context, p *ProductList) (int, error)
	GetAllNames(ctx context.Context, UserID int64) ([]ProductList, error)
	GetFoolInfoGroupProdList(ctx context.Context, listID int) (*FoolInfoProductList, error)
	GetFoolInfoProdList(ctx context.Context, listID int) (*FoolInfoProductList, error)
	MergeListGroup(ctx context.Context, listID, groupID int) error
	MakeListInactive(ctx context.Context, listID int) error
	MakeListActive(ctx context.Context, listID int) (string, error)
	Delete(ctx context.Context, listID int) error
}

type ProductRepo interface {
	Add(prod []Product) error
	CountByListID(c context.Context, listID int) (int, error)
	GetByListID(context context.Context, listID, offset, limit int) ([]Product, error)
	Delete(context context.Context, id int) error
}

type User struct {
	ChatID   int64
	UserName *string
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

type Product struct {
	ID      int
	Product string
	UserID  int64
	ListID  *int
}

type FoolInfoProductList struct {
	List     ProductList
	Products []InfoProduct
}

type InfoProduct struct {
	Product string
	User    *User
}

type Editors struct {
	User            User
	ManyAddProducts int
}

type GroupList struct {
	PruductLists []ProductList
	GroupOwnerID int64
}
