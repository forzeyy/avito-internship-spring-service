package handlers

import (
	"net/http"

	"github.com/forzeyy/avito-internship-spring-service/internal/dto"
	"github.com/forzeyy/avito-internship-spring-service/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	userSvc services.UserService
}

func NewAuthHandler(userSvc services.UserService) *AuthHandler {
	return &AuthHandler{
		userSvc: userSvc,
	}
}

// POST /register
func (ah *AuthHandler) RegisterUser(c echo.Context) error {
	var request dto.PostRegisterJSONRequestBody
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: "невалидный запрос",
		})
	}

	err := ah.userSvc.RegisterUser(c.Request().Context(), string(request.Email), request.Password, string(request.Role))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: err.Error(),
		})
	}
	return c.NoContent(http.StatusCreated)
}

// POST /login
func (ah *AuthHandler) LoginUser(c echo.Context) error {
	var request dto.PostLoginJSONRequestBody
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: "невалидный запрос",
		})
	}

	token, err := ah.userSvc.LoginUser(c.Request().Context(), string(request.Email), request.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.Error{
			Message: "неверные учетные данные",
		})
	}
	return c.JSON(http.StatusOK, dto.Token(token))
}

// POST /dummyLogin
func (ah *AuthHandler) DummyLogin(c echo.Context) error {
	var request dto.PostDummyLoginJSONRequestBody
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: "невалидный запрос",
		})
	}

	token, err := ah.userSvc.DummyLogin(c.Request().Context(), string(request.Role))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.Token(token))
}
