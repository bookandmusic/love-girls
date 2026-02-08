package error

import "errors"

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token expired")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEncrypt            = errors.New("encrypt error")
	ErrDecrypt            = errors.New("decrypt error")
)
