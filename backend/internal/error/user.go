package error

import "errors"

var (
	ErrUserExists  = errors.New("user already exists")
	ErrEmailExists = errors.New("email already exists")
)
