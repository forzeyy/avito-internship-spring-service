package utils

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GenerateAccessToken(role, secret string) (string, error) {
	if len(secret) == 0 {
		return "", errors.New("секрет JWT не может быть пустым")
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	return accessToken.SignedString([]byte(secret))
}

func VerifyAccessToken(tokenString string, secret string) (*string, error) {
	if len(secret) == 0 {
		return nil, errors.New("секрет JWT не может быть пустым")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role, ok := claims["role"].(string)
		if !ok || role == "" {
			return nil, jwt.ErrInvalidKey
		}

		return &role, nil
	}

	return nil, jwt.ErrTokenMalformed
}

func GetRoleFromContext(c echo.Context) (string, error) {
	token := c.Get("user").(*jwt.Token)
	if token == nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Неавторизован.")
	}

	claims := token.Claims.(jwt.MapClaims)
	role, ok := claims["role"].(string)
	if !ok {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Неавторизован.")
	}

	return role, nil
}
