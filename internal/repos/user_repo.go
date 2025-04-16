package repos

import (
	"context"
	"errors"
	"fmt"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepo struct {
	db DB
}

func NewUserRepo(db DB) UserRepo {
	return &userRepo{db: db}
}

func (ur *userRepo) CreateUser(ctx context.Context, user *models.User) error {
	if user.PasswordHash == "" {
		return errors.New("пароль пользователя не может быть пустым")
	}

	query := `
		INSERT INTO users (email, password, role, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := ur.db.Exec(ctx, query, user.Email, user.PasswordHash, user.Role, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("не удалось создать пользователя: %v", err)
	}
	return nil
}

func (ur *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := `
        SELECT id, email, password_hash, role, created_at
        FROM users
        WHERE email = $1
    `
	row := ur.db.QueryRow(ctx, query, email)

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, errors.New("пользователь не найден")
	}
	if err != nil {
		return nil, fmt.Errorf("не удалось получить пользователя: %v", err)
	}
	return &user, nil
}
