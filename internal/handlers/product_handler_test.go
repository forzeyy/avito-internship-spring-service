package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/forzeyy/avito-internship-spring-service/internal/handlers"
	"github.com/forzeyy/avito-internship-spring-service/internal/mocks"
	"github.com/forzeyy/avito-internship-spring-service/internal/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddProduct_Success(t *testing.T) {
	e := echo.New()
	mockService := new(mocks.ProductService)
	handler := handlers.NewProductHandler(mockService)

	pvzID := uuid.New()
	receptionID := uuid.New()
	date := time.Now()

	payload := `{"type":"одежда", "pvzId":"` + pvzID.String() + `"}`

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	mockService.On("AddProduct", mock.Anything, "одежда", pvzID.String()).Return(&models.Product{
		ID:          uuid.New(),
		Type:        "одежда",
		ReceptionID: receptionID,
		DateTime:    date,
	}, nil)

	err := handler.AddProduct(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestAddProduct_BindFail(t *testing.T) {
	e := echo.New()
	mockService := new(mocks.ProductService)
	handler := handlers.NewProductHandler(mockService)

	payload := `{"type":123}`

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler.AddProduct(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "невалидный запрос")
}

func TestAddProduct_ServiceError(t *testing.T) {
	e := echo.New()
	mockService := new(mocks.ProductService)
	handler := handlers.NewProductHandler(mockService)

	pvzID := uuid.New()
	payload := `{"type":"обувь", "pvzId":"` + pvzID.String() + `"}`

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	mockService.On("AddProduct", mock.Anything, "обувь", pvzID.String()).Return(&models.Product{}, errors.New("ошибка добавления"))

	err := handler.AddProduct(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "ошибка добавления")
}

func TestDeleteLastProduct_Success(t *testing.T) {
	e := echo.New()
	mockService := new(mocks.ProductService)
	handler := handlers.NewProductHandler(mockService)

	pvzID := uuid.New().String()

	req := httptest.NewRequest(http.MethodDelete, "/pvz/"+pvzID+"/product", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("pvzId")
	ctx.SetParamValues(pvzID)

	mockService.On("DeleteLastProduct", mock.Anything, pvzID).Return(nil)

	err := handler.DeleteLastProduct(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteLastProduct_InvalidUUID(t *testing.T) {
	e := echo.New()
	mockService := new(mocks.ProductService)
	handler := handlers.NewProductHandler(mockService)

	invalidID := "not-a-uuid"

	req := httptest.NewRequest(http.MethodDelete, "/pvz/"+invalidID+"/product", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("pvzId")
	ctx.SetParamValues(invalidID)

	err := handler.DeleteLastProduct(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "неверный формат pvz_id")
}

func TestDeleteLastProduct_ServiceError(t *testing.T) {
	e := echo.New()
	mockService := new(mocks.ProductService)
	handler := handlers.NewProductHandler(mockService)

	pvzID := uuid.New().String()

	req := httptest.NewRequest(http.MethodDelete, "/pvz/"+pvzID+"/product", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("pvzId")
	ctx.SetParamValues(pvzID)

	mockService.On("DeleteLastProduct", mock.Anything, pvzID).Return(errors.New("ошибка удаления"))

	err := handler.DeleteLastProduct(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "ошибка удаления")
}
