package store

import "errors"

var (
	NoRowProductError       = errors.New("List not exists")
	NoRowListOfProductError = errors.New("users lists is empty")
)
