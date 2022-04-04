package auth

import (
	"fmt"

	"github.com/VanjaRo/web-chat/config"
	"github.com/VanjaRo/web-chat/interfaces"
	"github.com/golang-jwt/jwt"
)

// 1 week expire time
const defaulExpireTime = 604800

// config.GoDotEnvVariable("JWT_SECRET")

type Claims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

func (claims *Claims) GetId() string {
	return claims.ID
}

func (claims *Claims) GetName() string {
	return claims.Name
}

func GenerateJWTToken(user interfaces.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":        user.GetId(),
		"Name":      user.GetName(),
		"ExpiresAt": jwt.TimeFunc().Add(defaulExpireTime).Unix(),
	})
	// tokenString, err
	return token.SignedString([]byte(config.GoDotEnvVariable("JWT_SECRET")))
}

func ValidateJWTToken(tokenString string) (interfaces.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// check the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GoDotEnvVariable("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
