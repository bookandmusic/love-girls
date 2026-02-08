package error

import "errors"

var (
	ErrQuery    = errors.New("query failed")
	ErrCreate   = errors.New("create failed")
	ErrUpdate   = errors.New("update failed")
	ErrDelete   = errors.New("delete failed")
	ErrNotFound = errors.New("not found")
)
