package structs

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
