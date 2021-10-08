package handlers

import (
	"github.com/golang-jwt/jwt"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type JwtCustomClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type Handler struct {
	db *gorm.DB
	ec *nats.EncodedConn
}

func New(db *gorm.DB, ec *nats.EncodedConn) *Handler {
	return &Handler{db: db, ec: ec}
}
