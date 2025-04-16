package services

import (
	"context"
	"errors"
	"time"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/forzeyy/avito-internship-spring-service/internal/repos"
	"github.com/google/uuid"
)

type ReceptionService interface {
	CreateReception(ctx context.Context, pvzID string) (*models.Reception, error)
	CloseLastReception(ctx context.Context, pvzID string) (*models.Reception, error)
}

type receptionService struct {
	receptionRepo repos.ReceptionRepo
}

func NewReceptionService(receptionRepo repos.ReceptionRepo) ReceptionService {
	return &receptionService{receptionRepo: receptionRepo}
}

func (rs *receptionService) CreateReception(ctx context.Context, pvzID string) (*models.Reception, error) {
	_, err := uuid.Parse(pvzID)
	if err != nil {
		return nil, errors.New("неверный формат pvz_id")
	}

	reception := &models.Reception{
		ID:       uuid.New(),
		PVZID:    uuid.MustParse(pvzID),
		Status:   "in_progress",
		DateTime: time.Now(),
	}

	err = rs.receptionRepo.CreateReception(ctx, reception)
	if err != nil {
		return nil, err
	}

	return reception, nil
}

func (rs *receptionService) CloseLastReception(ctx context.Context, pvzID string) (*models.Reception, error) {
	_, err := uuid.Parse(pvzID)
	if err != nil {
		return nil, errors.New("неверный формат pvz_id")
	}

	return rs.receptionRepo.CloseLastReception(ctx, uuid.MustParse(pvzID))
}
