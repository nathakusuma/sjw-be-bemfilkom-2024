package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type JWT struct {
	SecretKey []byte
	TTL       time.Duration
}

type Claims struct {
	jwt.RegisteredClaims
	Role         string `json:"role"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	ProgramStudi string `json:"program_studi"`
	Angkatan     string `json:"angkatan"`
}

func NewJWT(secretKey string, ttlString string) JWT {
	ttl, err := time.ParseDuration(ttlString)
	if err != nil || ttl <= 0 {
		log.Fatalln(err)
	}

	return JWT{
		SecretKey: []byte(secretKey),
		TTL:       ttl,
	}
}

func (j *JWT) Create(claims *Claims) (string, error) {
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(j.TTL))

	unsignedJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := unsignedJWT.SignedString(j.SecretKey)
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func (j *JWT) Decode(tokenString string, claims *Claims) error {
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
