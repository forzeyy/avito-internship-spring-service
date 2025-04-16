package services

import (
	"context"
	"errors"
	"time"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/forzeyy/avito-internship-spring-service/internal/repos"
	"github.com/google/uuid"
)

type PVZService interface {
	CreatePVZ(ctx context.Context, city string) (*models.PVZ, error)
	GetPVZs(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error)
}

type pvzService struct {
	pvzRepo repos.PVZRepo
}

func NewPVZService(pvzRepo repos.PVZRepo) PVZService {
	return &pvzService{pvzRepo: pvzRepo}
}

func (ps *pvzService) CreatePVZ(ctx context.Context, city string) (*models.PVZ, error) {
	if city != "Москва" && city != "Санкт-Петербург" && city != "Казань" {
		return nil, errors.New("неверный город")
	}

	pvz := &models.PVZ{
		ID:      uuid.New(),
		City:    city,
		RegDate: time.Now(),
	}

	err := ps.pvzRepo.CreatePVZ(ctx, pvz)
	if err != nil {
		return nil, err
	}

	return pvz, nil
}

func (ps *pvzService) GetPVZs(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error) {
	return ps.pvzRepo.GetPVZs(ctx, startDate, endDate, page, limit)
}
