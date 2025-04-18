package repos

import (
	"context"
	"errors"
	"fmt"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/google/uuid"
)

type ProductRepo interface {
	AddProduct(ctx context.Context, product *models.Product) error
	DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error
}

type productRepo struct {
	db DB
}

func NewProductRepo(db DB) ProductRepo {
	return &productRepo{db: db}
}

func (pr *productRepo) AddProduct(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (id, type, reception_id, received_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := pr.db.Exec(ctx, query, product.ID, product.Type, product.ReceptionID, product.DateTime)
	if err != nil {
		return fmt.Errorf("не удалось добавить продукт: %v", err)
	}
	return nil
}

func (pr *productRepo) DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error {
	query := `
        WITH last_product AS (
            SELECT id
            FROM products
            WHERE reception_id = (
                SELECT id
                FROM receptions
                WHERE pvz_id = $1 AND status = 'in_progress'
                ORDER BY created_at DESC
                LIMIT 1
            )
            ORDER BY received_at DESC
            LIMIT 1
        )
        DELETE FROM products
        WHERE id = (SELECT id FROM last_product)
    `

	result, err := pr.db.Exec(ctx, query, pvzID)
	if err != nil {
		return fmt.Errorf("не удалось удалить последний товар: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("нет товаров для удаления")
	}
	return nil
}
