package jwt

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/entity"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type IJWT interface {
	Create(user entity.User) (string, error)
	Decode(tokenString string, claims *Claims) error
}

type JWT struct {
	SecretKey []byte
	TTL       time.Duration
}

type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

func NewJWT(secretKey string, ttlString string) IJWT {
	ttl, err := time.ParseDuration(ttlString)
	if err != nil || ttl <= 0 {
		log.Fatalln(err)
	}

	return JWT{
		SecretKey: []byte(secretKey),
		TTL:       ttl,
	}
}

func (j JWT) Create(user entity.User) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Nim,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.TTL)),
		},
		Role: user.Role,
	}

	unsignedJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := unsignedJWT.SignedString(j.SecretKey)
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func (j JWT) Decode(tokenString string, claims *Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.SecretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}
