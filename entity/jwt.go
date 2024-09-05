package entity

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	TokenType string `json:"token_type"`
	UserID    int64  `json:"user_id"`
}
