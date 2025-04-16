package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/forzeyy/avito-internship-spring-service/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPVZRepo struct {
	mock.Mock
}

func (m *mockPVZRepo) CreatePVZ(ctx context.Context, pvz *models.PVZ) error {
	args := m.Called(ctx, pvz)
	return args.Error(0)
}

func (m *mockPVZRepo) GetPVZs(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error) {
	args := m.Called(ctx, startDate, endDate, page, limit)
	return args.Get(0).([]models.PVZ), args.Error(1)
}

// CreatePVZ
func TestCreatePVZ_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockPVZRepo)
	service := services.NewPVZService(mockRepo)

	city := "Москва"

	mockRepo.On("CreatePVZ", ctx, mock.AnythingOfType("*models.PVZ")).Return(nil)

	pvz, err := service.CreatePVZ(ctx, city)

	assert.NoError(t, err)
	assert.NotNil(t, pvz)
	assert.Equal(t, city, pvz.City)
	mockRepo.AssertExpectations(t)
}

func TestCreatePVZ_InvalidCity(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockPVZRepo)
	service := services.NewPVZService(mockRepo)

	pvz, err := service.CreatePVZ(ctx, "Париж")

	assert.Nil(t, pvz)
	assert.EqualError(t, err, "неверный город")
	mockRepo.AssertNotCalled(t, "CreatePVZ")
}

func TestCreatePVZ_RepoError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockPVZRepo)
	service := services.NewPVZService(mockRepo)

	mockRepo.On("CreatePVZ", ctx, mock.AnythingOfType("*models.PVZ")).Return(errors.New("db error"))

	pvz, err := service.CreatePVZ(ctx, "Казань")

	assert.Nil(t, pvz)
	assert.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}

// GetPVZs
func TestGetPVZs_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockPVZRepo)
	service := services.NewPVZService(mockRepo)

	now := time.Now()
	expected := []models.PVZ{
		{ID: uuid.New(), City: "Москва", RegDate: now},
		{ID: uuid.New(), City: "Казань", RegDate: now},
	}

	mockRepo.On("GetPVZs", ctx, &now, &now, 1, 10).Return(expected, nil)

	result, err := service.GetPVZs(ctx, &now, &now, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}
