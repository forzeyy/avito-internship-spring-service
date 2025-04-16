package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/forzeyy/avito-internship-spring-service/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReceptionRepo struct {
	mock.Mock
}

func (m *mockReceptionRepo) CreateReception(ctx context.Context, reception *models.Reception) error {
	args := m.Called(ctx, reception)
	return args.Error(0)
}

func (m *mockReceptionRepo) GetLastOpenReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	args := m.Called(ctx, pvzID)
	if rec, ok := args.Get(0).(*models.Reception); ok {
		return rec, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockReceptionRepo) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	args := m.Called(ctx, pvzID)
	if rec, ok := args.Get(0).(*models.Reception); ok {
		return rec, args.Error(1)
	}
	return nil, args.Error(1)
}

// CreateReception
func TestCreateReception_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockReceptionRepo)
	service := services.NewReceptionService(mockRepo)

	pvzID := uuid.New().String()

	mockRepo.On("CreateReception", ctx, mock.AnythingOfType("*models.Reception")).Return(nil)

	reception, err := service.CreateReception(ctx, pvzID)

	assert.NoError(t, err)
	assert.NotNil(t, reception)
	assert.Equal(t, "in_progress", reception.Status)
	mockRepo.AssertExpectations(t)
}

func TestCreateReception_InvalidUUID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockReceptionRepo)
	service := services.NewReceptionService(mockRepo)

	reception, err := service.CreateReception(ctx, "invalid-uuid")

	assert.Nil(t, reception)
	assert.EqualError(t, err, "неверный формат pvz_id")
	mockRepo.AssertNotCalled(t, "CreateReception")
}

func TestCreateReception_RepoError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockReceptionRepo)
	service := services.NewReceptionService(mockRepo)

	pvzID := uuid.New().String()

	mockRepo.On("CreateReception", ctx, mock.Anything).Return(errors.New("db error"))

	reception, err := service.CreateReception(ctx, pvzID)

	assert.Nil(t, reception)
	assert.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}

// CloseLastReception
func TestCloseLastReception_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockReceptionRepo)
	service := services.NewReceptionService(mockRepo)

	pvzID := uuid.New()
	expectedReception := &models.Reception{
		ID:     uuid.New(),
		PVZID:  pvzID,
		Status: "closed",
	}

	mockRepo.On("CloseLastReception", ctx, pvzID).Return(expectedReception, nil)

	reception, err := service.CloseLastReception(ctx, pvzID.String())

	assert.NoError(t, err)
	assert.Equal(t, expectedReception, reception)
	mockRepo.AssertExpectations(t)
}

func TestCloseLastReception_InvalidUUID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockReceptionRepo)
	service := services.NewReceptionService(mockRepo)

	reception, err := service.CloseLastReception(ctx, "not-a-uuid")

	assert.Nil(t, reception)
	assert.EqualError(t, err, "неверный формат pvz_id")
	mockRepo.AssertNotCalled(t, "CloseLastReception")
}

func TestCloseLastReception_RepoError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockReceptionRepo)
	service := services.NewReceptionService(mockRepo)

	pvzID := uuid.New()
	mockRepo.On("CloseLastReception", ctx, pvzID).Return(nil, errors.New("close error"))

	reception, err := service.CloseLastReception(ctx, pvzID.String())

	assert.Nil(t, reception)
	assert.EqualError(t, err, "close error")
	mockRepo.AssertExpectations(t)
}
