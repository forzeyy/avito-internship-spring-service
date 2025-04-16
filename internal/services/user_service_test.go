package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/forzeyy/avito-internship-spring-service/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

type mockAuthUtil struct {
	mock.Mock
}

func (m *mockAuthUtil) GenerateAccessToken(role, key string) (string, error) {
	args := m.Called(role, key)
	return args.String(0), args.Error(1)
}

func (m *mockAuthUtil) CheckPassword(hashedPassword, password string) bool {
	args := m.Called(hashedPassword, password)
	return args.Bool(0)
}

// RegisterUser
func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockAuth := new(mockAuthUtil)

	email := "defaultuser@mail.ru"
	password := "WeakPassword123"
	role := "employee"

	mockRepo.On("GetUserByEmail", mock.Anything, email).Return(nil, errors.New("not found"))
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	svc := services.NewUserService(mockRepo, "secret", mockAuth)

	err := svc.RegisterUser(context.Background(), email, password, role)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockAuth := new(mockAuthUtil)

	email := "exists@mail.ru"

	mockRepo.On("GetUserByEmail", mock.Anything, email).Return(&models.User{}, nil)

	svc := services.NewUserService(mockRepo, "secret", mockAuth)

	err := svc.RegisterUser(context.Background(), email, "passwd", "employee")
	assert.EqualError(t, err, "пользователь с таким email уже существует")
}

func TestRegisterUser_InvalidRole(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockAuth := new(mockAuthUtil)

	svc := services.NewUserService(mockRepo, "secret", mockAuth)

	err := svc.RegisterUser(context.Background(), "admin@mail.ru", "123", "admin")
	assert.EqualError(t, err, "неверная роль пользователя")
}

// LoginUser
func TestLoginUser_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockAuth := new(mockAuthUtil)

	email := "test@mail.ru"
	password := "password"
	hashed := "hashedpassword"
	role := "moderator"

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashed,
		Role:         role,
		CreatedAt:    time.Now(),
	}

	mockRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
	mockAuth.On("CheckPassword", hashed, password).Return(true)
	mockAuth.On("GenerateAccessToken", role, "secret").Return("token123", nil)

	svc := services.NewUserService(mockRepo, "secret", mockAuth)

	token, err := svc.LoginUser(context.Background(), email, password)
	assert.NoError(t, err)
	assert.Equal(t, "token123", token)

	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockAuth := new(mockAuthUtil)

	email := "forgetful.user@mail.ru"
	password := "wrong"
	hashed := "hashed"

	user := &models.User{PasswordHash: hashed}

	mockRepo.On("GetUserByEmail", mock.Anything, email).Return(user, nil)
	mockAuth.On("CheckPassword", hashed, password).Return(false)

	svc := services.NewUserService(mockRepo, "secret", mockAuth)

	token, err := svc.LoginUser(context.Background(), email, password)
	assert.Empty(t, token)
	assert.EqualError(t, err, "неверный пароль")
}

func TestLoginUser_NotFound(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockAuth := new(mockAuthUtil)

	mockRepo.On("GetUserByEmail", mock.Anything, "user404@mail.ru").Return(nil, errors.New("not found"))

	svc := services.NewUserService(mockRepo, "secret", mockAuth)

	token, err := svc.LoginUser(context.Background(), "user404@mail.ru", "any")
	assert.Empty(t, token)
	assert.EqualError(t, err, "пользователь с таким email не найден")
}

// DummyLogin
func TestDummyLogin_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockAuth := new(mockAuthUtil)

	mockAuth.On("GenerateAccessToken", "employee", "secret").Return("dummy_token", nil)

	svc := services.NewUserService(mockRepo, "secret", mockAuth)

	token, err := svc.DummyLogin(context.Background(), "employee")
	assert.NoError(t, err)
	assert.Equal(t, "dummy_token", token)
}

func TestDummyLogin_InvalidRole(t *testing.T) {
	mockRepo := new(mockUserRepo)
	mockAuth := new(mockAuthUtil)

	svc := services.NewUserService(mockRepo, "secret", mockAuth)

	token, err := svc.DummyLogin(context.Background(), "admin")
	assert.Empty(t, token)
	assert.EqualError(t, err, "неверная роль пользователя")
}
