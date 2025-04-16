package repos_test

import (
	"context"
	"testing"
	"time"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/forzeyy/avito-internship-spring-service/internal/repos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

// CreateUser
func TestCreateUser_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewUserRepo(mock)

	user := &models.User{
		ID:           uuid.New(),
		Email:        "test@mail.ru",
		PasswordHash: "passhash123123",
		Role:         "employee",
		CreatedAt:    time.Now(),
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Email, user.PasswordHash, user.Role, user.CreatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_MissingPassword(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := repos.NewUserRepo(mock)

	user := &models.User{
		ID:        uuid.New(),
		Email:     "test@mail.ru",
		Role:      "employee",
		CreatedAt: time.Now(),
	}

	err := repo.CreateUser(context.Background(), user)
	assert.EqualError(t, err, "пароль пользователя не может быть пустым")
}

// GetUserByEmail
func TestGetUserByEmail_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repos.NewUserRepo(mock)

	user := &models.User{
		ID:           uuid.New(),
		Email:        "test@mail.ru",
		PasswordHash: "passhash123123",
		Role:         "employee",
		CreatedAt:    time.Now(),
	}

	rows := pgxmock.NewRows([]string{"id", "email", "password_hash", "role", "created_at"}).
		AddRow(user.ID, user.Email, user.PasswordHash, user.Role, user.CreatedAt)

	mock.ExpectQuery("SELECT id, email, password_hash, role, created_at FROM users").
		WithArgs(user.Email).
		WillReturnRows(rows)

	result, err := repo.GetUserByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, result.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := repos.NewUserRepo(mock)

	mock.ExpectQuery("SELECT id, email, password_hash, role, created_at FROM users").
		WithArgs("404@mail.ru").
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.GetUserByEmail(context.Background(), "404@mail.ru")
	assert.Nil(t, result)
	assert.EqualError(t, err, "пользователь не найден")
}
