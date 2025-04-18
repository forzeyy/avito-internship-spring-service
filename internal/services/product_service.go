package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/forzeyy/avito-internship-spring-service/internal/repos"
	"github.com/google/uuid"
)

type ProductService interface {
	AddProduct(ctx context.Context, productType, pvzID string) (*models.Product, error)
	DeleteLastProduct(ctx context.Context, pvzID string) error
}

type productService struct {
	prodRepo repos.ProductRepo
	recRepo  repos.ReceptionRepo
}

func NewProductService(prodRepo repos.ProductRepo, recRepo repos.ReceptionRepo) ProductService {
	return &productService{
		prodRepo: prodRepo,
		recRepo:  recRepo,
	}
}

func (ps *productService) AddProduct(ctx context.Context, productType, pvzID string) (*models.Product, error) {
	if productType != "электроника" && productType != "одежда" && productType != "обувь" {
		return nil, errors.New("недопустимый тип товара")
	}

	parsedPVZID, err := uuid.Parse(pvzID)
	if err != nil || parsedPVZID == uuid.Nil {
		return nil, errors.New("неверный формат pvz_id")
	}

	lastReception, err := ps.recRepo.GetLastOpenReception(ctx, parsedPVZID)
	if err != nil {
		return nil, err
	}
	if lastReception == nil {
		return nil, errors.New("последняя открытая приемка не найдена")
	}

	product := &models.Product{
		ID:          uuid.New(),
		Type:        productType,
		ReceptionID: lastReception.ID,
		DateTime:    time.Now(),
	}
	err = ps.prodRepo.AddProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (ps *productService) DeleteLastProduct(ctx context.Context, pvzID string) error {
	parsedPVZID, err := uuid.Parse(pvzID)
	if err != nil {
		return errors.New("неверный формат pvz_id")
	}

	// проверка есть ли открытая приемка
	lastReception, err := ps.recRepo.GetLastOpenReception(ctx, parsedPVZID)
	if err != nil {
		return fmt.Errorf("не удалось получить последнюю открытую приемку: %w", err)
	}
	if lastReception == nil || lastReception.ID == uuid.Nil {
		return errors.New("последняя открытая приемка не найдена")
	}

	return ps.prodRepo.DeleteLastProduct(ctx, parsedPVZID)
}
