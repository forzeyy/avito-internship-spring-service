package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/forzeyy/avito-internship-spring-service/internal/handlers"
	"github.com/forzeyy/avito-internship-spring-service/internal/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser(t *testing.T) {
	e := echo.New()
	mockUserService := new(mocks.UserService)
	handler := handlers.NewAuthHandler(mockUserService)

	payload := `{"email":"user@example.com","password":"secret","role":"employee"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserService.On("RegisterUser", mock.Anything, "user@example.com", "secret", "employee").Return(nil)

	err := handler.RegisterUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUserService.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	e := echo.New()
	mockUserService := new(mocks.UserService)
	handler := handlers.NewAuthHandler(mockUserService)

	payload := `{"email":"user@example.com","password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserService.On("LoginUser", mock.Anything, "user@example.com", "secret").Return("mock_token", nil)

	err := handler.LoginUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `"mock_token"`, rec.Body.String())
	mockUserService.AssertExpectations(t)
}

func TestLoginUser_Fail(t *testing.T) {
	e := echo.New()
	mockUserService := new(mocks.UserService)
	handler := handlers.NewAuthHandler(mockUserService)

	payload := `{"email":"user@example.com","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserService.On("LoginUser", mock.Anything, "user@example.com", "wrong").Return("", errors.New("invalid"))

	err := handler.LoginUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	mockUserService.AssertExpectations(t)
}

func TestDummyLogin(t *testing.T) {
	e := echo.New()
	mockUserService := new(mocks.UserService)
	handler := handlers.NewAuthHandler(mockUserService)

	payload := `{"role":"moderator"}`
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserService.On("DummyLogin", mock.Anything, "moderator").Return("dummy_token", nil)

	err := handler.DummyLogin(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `"dummy_token"`, rec.Body.String())
	mockUserService.AssertExpectations(t)
}
