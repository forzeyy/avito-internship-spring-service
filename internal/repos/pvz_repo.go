package repos

import (
	"context"
	"fmt"
	"time"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
)

type PVZRepo interface {
	CreatePVZ(ctx context.Context, pvz *models.PVZ) error
	GetPVZs(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error)
}

type pvzRepo struct {
	db DB
}

func NewPVZRepo(db DB) PVZRepo {
	return &pvzRepo{db: db}
}

func (pr *pvzRepo) CreatePVZ(ctx context.Context, pvz *models.PVZ) error {
	query := `
		INSERT INTO pvzs (id, city, reg_date)
		VALUES ($1, $2, $3)
	`
	_, err := pr.db.Exec(ctx, query, pvz.ID, pvz.City, pvz.RegDate)
	if err != nil {
		return fmt.Errorf("не удалось создать ПВЗ: %v", err)
	}
	return nil
}

func (pr *pvzRepo) GetPVZs(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error) {
	var (
		query = `
			SELECT id, city, reg_date
			FROM pvzs
			WHERE city IN ('Москва', 'Санкт-Петербург', 'Казань')
		`
		args     []any
		argIndex = 1
	)

	if startDate != nil {
		query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM receptions WHERE receptions.pvz_id = pvzs.id AND receptions.created_at >= $%d)", argIndex)
		args = append(args, startDate)
		argIndex++
	}

	if endDate != nil {
		query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM receptions WHERE receptions.pvz_id = pvzs.id AND receptions.created_at <= $%d)", argIndex)
		args = append(args, endDate)
		argIndex++
	}

	query += fmt.Sprintf(`
        ORDER BY reg_date DESC
        LIMIT $%d OFFSET $%d
    `, argIndex, argIndex+1)

	args = append(args, limit, (page-1)*limit)

	rows, err := pr.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка ПВЗ: %w", err)
	}
	defer rows.Close()

	var pvzs []models.PVZ
	for rows.Next() {
		var pvz models.PVZ
		err := rows.Scan(&pvz.ID, &pvz.City, &pvz.RegDate)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}
		pvzs = append(pvzs, pvz)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения строк: %v", err)
	}
	return pvzs, nil
}
