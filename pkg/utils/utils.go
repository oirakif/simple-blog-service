package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("mysecretkey123") // The secret key used for signing the token

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTUtils struct {
	SecretKey []byte
}

func NewJWTUtils(secretKey []byte) *JWTUtils {

	return &JWTUtils{
		SecretKey: secretKey,
	}
}

func (jwtu *JWTUtils) GenerateJWT(email string, claims jwt.Claims) (signedToken string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

// Function to parse and validate JWT
func (jwtu *JWTUtils) ValidateToken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	splittedToken := strings.Split(token, "Bearer ")
	if len(splittedToken) != 2 {
		c.Status(http.StatusUnauthorized)
		c.Abort()
		return
	}
	// Parse the token
	parsed, err := jwt.ParseWithClaims(splittedToken[1], &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for validation
		return secretKey, nil
	})

	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	// Validate token and check the claims
	if claims, ok := parsed.Claims.(*JWTClaims); ok && parsed.Valid {
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()

		return
	}
	c.String(http.StatusUnauthorized, "invalid token")
	c.Abort()
}

func HashSHA256(plaintext string) (hashed string) {
	h := sha256.New()
	h.Write([]byte(plaintext))
	hashed = hex.EncodeToString(h.Sum(nil))

	return hashed
}

func CalculateOffset(page, perPage int) int {
	if page < 1 || perPage < 1 {
		return 0
	}
	return (page - 1) * perPage
}
