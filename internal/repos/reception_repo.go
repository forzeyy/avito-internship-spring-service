package repos

import (
	"context"
	"errors"
	"fmt"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ReceptionRepo interface {
	CreateReception(ctx context.Context, reception *models.Reception) error
	CloseLastReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error)
}

type receptionRepo struct {
	db DB
}

func NewReceiptRepo(db DB) ReceptionRepo {
	return &receptionRepo{db: db}
}

func (rr *receptionRepo) CreateReception(ctx context.Context, reception *models.Reception) error {
	query := `
		INSERT INTO receptions (id, pvz_id, status, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := rr.db.Exec(ctx, query, reception.ID, reception.PVZID, reception.Status, reception.DateTime)
	if err != nil {
		return fmt.Errorf("не удалось создать приемку: %v", err)
	}
	return nil
}

func (rr *receptionRepo) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	var reception models.Reception

	query := `
		WITH last_reception AS (
			SELECT id
			FROM receipts
			WHERE pvz_id = $1 and status = 'in_progress'
			ORDER BY created_at DESC
			LIMIT 1
		)
		UPDATE receipts
		SET status = 'close', closed_at = NOW()
		WHERE id = (SELECT id FROM last_reception)
		RETURNING id, pvz_id, status, created_at, closed_at
	`
	err := rr.db.QueryRow(ctx, query, pvzID).Scan(
		&reception.ID,
		&reception.PVZID,
		&reception.Status,
		&reception.DateTime,
		&reception.ClosedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("не удалось закрыть последнюю приемку: %v", err)
	}
	return &reception, nil
}
