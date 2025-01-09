package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret   []byte
	tokenExpire time.Duration
)

// Initialize 初始化认证配置
func Initialize(secret string, expire time.Duration) {
	jwtSecret = []byte(secret)
	tokenExpire = expire
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userID uint) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenExpire).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   fmt.Sprint(userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		userID, _ := strconv.ParseUint(claims.Subject, 10, 32)
		return uint(userID), nil
	}

	return 0, jwt.ErrSignatureInvalid
}
