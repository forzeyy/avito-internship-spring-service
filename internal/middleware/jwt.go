package middleware

import (
	"log"
	"net/http"

	jwtMiddleware "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return jwtMiddleware.WithConfig(jwtMiddleware.Config{
		SigningKey:    []byte(jwtSecret),
		TokenLookup:   "header:Authorization:Bearer ",
		SigningMethod: "HS256",
		ErrorHandler: func(c echo.Context, err error) error {
			log.Printf("Ошибка проверки JWT: %v", err)
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid or expired jwt"})
		},
	})
}
