package ds

import (
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/google/uuid"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID     uint `json:"user_id"`
	IsLinguist bool `json:"is_linguist"`
}
