package auth

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	errMsg "github.com/bookandmusic/love-girl/internal/error"
)

type JWT interface {
	Generate(claims *Claims) (string, error)
	Parse(token string) (*Claims, error)
}

type HS256JWT struct {
	secret []byte
	issuer string
	expire int64
}

func NewHS256JWT(secret, issuer string, expire int64) JWT {
	return &HS256JWT{
		secret: []byte(secret),
		issuer: issuer,
		expire: expire,
	}
}

func (j *HS256JWT) Generate(claims *Claims) (string, error) {
	now := time.Now()
	ttl := time.Duration(j.expire) * time.Second
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": claims.UserID,
		"rol": claims.Role,
		"iss": j.issuer,
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
	})
	return token.SignedString(j.secret)
}

func (j *HS256JWT) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errMsg.ErrInvalidToken
		}
		return j.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errMsg.ErrExpiredToken
		}
		return nil, errMsg.ErrInvalidToken
	}

	if !token.Valid {
		return nil, errMsg.ErrInvalidToken
	}

	mc, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errMsg.ErrInvalidToken
	}

	uid, ok := mc["uid"].(float64)
	if !ok {
		return nil, errMsg.ErrInvalidToken
	}

	role, _ := mc["rol"].(string)

	return &Claims{
		UserID: uint64(uid),
		Role:   role,
	}, nil
}
