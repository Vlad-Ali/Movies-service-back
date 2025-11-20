package jwtclaims

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID   string `json:"userID"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
