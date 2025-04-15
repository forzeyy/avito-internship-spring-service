package repos

import (
	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) CreateUser(user *models.User) error {
	/*query = `
		INSERT INTO users (email, password, role, created_at)
		VALUES ($1, $2, $3, $4)
	` */
}
