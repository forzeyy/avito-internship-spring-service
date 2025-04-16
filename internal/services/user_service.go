package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/forzeyy/avito-internship-spring-service/internal/repos"
	"github.com/forzeyy/avito-internship-spring-service/internal/utils"
	"github.com/google/uuid"
)

type UserService interface {
	DummyLogin(ctx context.Context, role string) (string, error)
	RegisterUser(ctx context.Context, email, password, role string) error
	LoginUser(ctx context.Context, email, password string) (string, error)
}

type userService struct {
	userRepo repos.UserRepo
	jwtKey   []byte
	authUtil utils.AuthUtil
}

func NewUserService(userRepo repos.UserRepo, jwtKey string, authUtil utils.AuthUtil) UserService {
	return &userService{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
		authUtil: authUtil,
	}
}

func (us *userService) DummyLogin(ctx context.Context, role string) (string, error) {
	if role != "employee" && role != "moderator" {
		return "", errors.New("неверная роль пользователя")
	}

	token, err := us.authUtil.GenerateAccessToken(role, string(us.jwtKey))
	if err != nil {
		return "", fmt.Errorf("не удалось сгенерировать токен: %v", err)
	}

	return token, nil
}

func (us *userService) RegisterUser(ctx context.Context, email, password, role string) error {
	if role != "employee" && role != "moderator" {
		return errors.New("неверная роль пользователя")
	}

	_, err := us.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return errors.New("пользователь с таким email уже существует")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("не удалось захешировать пароль: %v", err)
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         role,
		CreatedAt:    time.Now(),
	}

	return us.userRepo.CreateUser(ctx, user)
}

func (us *userService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := us.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("пользователь с таким email не найден")
	}

	if !us.authUtil.CheckPassword(user.PasswordHash, password) {
		return "", errors.New("неверный пароль")
	}

	token, err := us.authUtil.GenerateAccessToken(user.Role, string(us.jwtKey))
	if err != nil {
		return "", fmt.Errorf("не удалось сгенерировать токен: %v", err)
	}

	return token, nil
}
