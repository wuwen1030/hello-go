package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/wuwen/hello-go/internal/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	config *config.JWTConfig
}

func New(config *config.JWTConfig) *Auth {
	return &Auth{config: config}
}

func (a *Auth) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (a *Auth) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *Auth) GenerateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.config.ExpireTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.config.Secret))
}

// ... 其他方法
