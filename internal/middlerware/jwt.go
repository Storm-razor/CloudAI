package middlerware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wwwzy/CloudAI/config"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// ---------------------------
// @brief 传入userID,生成JWT Token
// ---------------------------
func GenerateToken(userID uint) (string, error) {
	cfg := config.AppConfigInstance.JWT
	expirationTime := time.Now().Add(time.Duration(cfg.ExpirationHours) * time.Hour)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ---------------------------
// @brief 传入token,返回解析出的Claims指针和可能的错误
// ---------------------------
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.AppConfigInstance.JWT

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
