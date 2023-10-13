package store

import "errors"

var (
	NoRowProductError       = errors.New("List not exists")
	NoRowListOfProductError = errors.New("You are have`t any list, create first!")
	NoUserGroupError        = errors.New("You are have`t any group, create or join in youre first group!")
)
