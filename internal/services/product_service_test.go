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

type mockProductRepo struct {
	mock.Mock
}

func (m *mockProductRepo) AddProduct(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *mockProductRepo) DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error {
	args := m.Called(ctx, pvzID)
	return args.Error(0)
}

// AddProduct
func TestAddProduct_Success(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockRec := new(mockReceptionRepo)

	pvzID := uuid.New()
	productType := "электроника"

	mockRec.On("GetLastOpenReception", mock.Anything, pvzID).Return(&models.Reception{}, nil)
	mockProd.On("AddProduct", mock.Anything, mock.AnythingOfType("*models.Product")).Return(nil)

	svc := services.NewProductService(mockProd, mockRec)

	product, err := svc.AddProduct(context.Background(), productType, pvzID.String())
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, productType, product.Type)
	assert.Equal(t, pvzID, product.ReceptionID)
}

func TestAddProduct_InvalidType(t *testing.T) {
	svc := services.NewProductService(nil, nil)

	product, err := svc.AddProduct(context.Background(), "еда", uuid.New().String())
	assert.Nil(t, product)
	assert.EqualError(t, err, "недопустимый тип товара")
}

func TestAddProduct_InvalidUUID(t *testing.T) {
	svc := services.NewProductService(nil, nil)

	product, err := svc.AddProduct(context.Background(), "одежда", "invalid-uuid")
	assert.Nil(t, product)
	assert.EqualError(t, err, "неверный формат pvz_id")
}

func TestAddProduct_NoOpenReception(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockRec := new(mockReceptionRepo)

	pvzID := uuid.New()
	mockRec.On("GetLastOpenReception", mock.Anything, pvzID).Return(nil, nil)

	svc := services.NewProductService(mockProd, mockRec)

	product, err := svc.AddProduct(context.Background(), "обувь", pvzID.String())
	mockRec.AssertCalled(t, "GetLastOpenReception", mock.Anything, pvzID)
	assert.Nil(t, product)
	assert.EqualError(t, err, "последняя открытая приемка не найдена")
}

func TestAddProduct_DBError(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockRec := new(mockReceptionRepo)

	pvzID := uuid.New()
	mockRec.On("GetLastOpenReception", mock.Anything, pvzID).Return(&models.Reception{}, nil)
	mockProd.On("AddProduct", mock.Anything, mock.Anything).Return(errors.New("db error"))

	svc := services.NewProductService(mockProd, mockRec)

	product, err := svc.AddProduct(context.Background(), "электроника", pvzID.String())
	assert.Nil(t, product)
	assert.EqualError(t, err, "db error")
}

// DeleteLastProduct
func TestDeleteLastProduct_Success(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockRec := new(mockReceptionRepo)

	pvzID := uuid.New()
	mockRec.On("GetLastOpenReception", mock.Anything, pvzID).Return(&models.Reception{}, nil)
	mockProd.On("DeleteLastProduct", mock.Anything, pvzID).Return(nil)

	svc := services.NewProductService(mockProd, mockRec)

	err := svc.DeleteLastProduct(context.Background(), pvzID.String())
	assert.NoError(t, err)
}

func TestDeleteLastProduct_InvalidUUID(t *testing.T) {
	svc := services.NewProductService(nil, nil)

	err := svc.DeleteLastProduct(context.Background(), "invalid-uuid")
	assert.EqualError(t, err, "неверный формат pvz_id")
}

func TestDeleteLastProduct_NoOpenReception(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockRec := new(mockReceptionRepo)

	pvzID := uuid.New()
	mockRec.On("GetLastOpenReception", mock.Anything, pvzID).Return(nil, nil)

	svc := services.NewProductService(mockProd, mockRec)

	err := svc.DeleteLastProduct(context.Background(), pvzID.String())
	assert.EqualError(t, err, "приемка закрыта")
}

func TestDeleteLastProduct_DBError(t *testing.T) {
	mockProd := new(mockProductRepo)
	mockRec := new(mockReceptionRepo)

	pvzID := uuid.New()
	mockRec.On("GetLastOpenReception", mock.Anything, pvzID).Return(&models.Reception{}, nil)
	mockProd.On("DeleteLastProduct", mock.Anything, pvzID).Return(errors.New("delete error"))

	svc := services.NewProductService(mockProd, mockRec)

	err := svc.DeleteLastProduct(context.Background(), pvzID.String())
	assert.EqualError(t, err, "delete error")
}
