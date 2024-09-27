package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type JWTUtils struct {
	SecretKey string
}

func NewJWTUtils(secretKey string) (jwtUtils JWTUtils) {
	jwtUtils.SecretKey = secretKey
	return jwtUtils
}

func (jwtu *JWTUtils) GenerateJWT(email string, claims jwt.Claims) (signedToken string, err error) {
	// Create the token using the signing method and claims
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// Sign the token with the secret key
	signedToken, err = token.SignedString([]byte(jwtu.SecretKey))
	if err != nil {
		log.Println(err)
	}
	return signedToken, err
}

func HashSHA256(plaintext string) (hashed string) {
	h := sha256.New()
	h.Write([]byte(plaintext))
	hashed = hex.EncodeToString(h.Sum(nil))

	return hashed
}
