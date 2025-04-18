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

func TestCreatePVZ_Success(t *testing.T) {
	e := echo.New()
	mockSvc := new(mocks.PVZService)
	handler := handlers.NewPVZHandler(mockSvc)

	payload := `{"city":"Москва"}`
	regDate := time.Now()
	pvzID := uuid.New()

	mockSvc.On("CreatePVZ", mock.Anything, "Москва").Return(&models.PVZ{
		ID:      pvzID,
		City:    "Москва",
		RegDate: regDate,
	}, nil)

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler.CreatePVZ(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"city":"Москва"`)
	mockSvc.AssertExpectations(t)
}

func TestCreatePVZ_BindError(t *testing.T) {
	e := echo.New()
	mockSvc := new(mocks.PVZService)
	handler := handlers.NewPVZHandler(mockSvc)

	payload := `{"city":123}` // невалидный формат

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler.CreatePVZ(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "невалидный запрос")
}

func TestCreatePVZ_ServiceError(t *testing.T) {
	e := echo.New()
	mockSvc := new(mocks.PVZService)
	handler := handlers.NewPVZHandler(mockSvc)

	payload := `{"city":"Питер"}`

	mockSvc.On("CreatePVZ", mock.Anything, "Питер").Return(&models.PVZ{}, errors.New("ошибка создания"))

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler.CreatePVZ(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "ошибка создания")
}

func TestGetPVZs_Success(t *testing.T) {
	e := echo.New()
	mockSvc := new(mocks.PVZService)
	handler := handlers.NewPVZHandler(mockSvc)

	regDate := time.Now()
	pvzID := uuid.New()
	pvzs := []models.PVZ{
		{
			ID:      pvzID,
			City:    "Москва",
			RegDate: regDate,
		},
	}

	mockSvc.On("GetPVZs", mock.Anything, (*time.Time)(nil), (*time.Time)(nil), 1, 10).Return(pvzs, nil)

	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler.GetPVZs(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"city":"Москва"`)
	mockSvc.AssertExpectations(t)
}

func TestGetPVZs_ServiceError(t *testing.T) {
	e := echo.New()
	mockSvc := new(mocks.PVZService)
	handler := handlers.NewPVZHandler(mockSvc)

	mockSvc.On("GetPVZs", mock.Anything, (*time.Time)(nil), (*time.Time)(nil), 1, 10).Return(nil, errors.New("ошибка получения"))

	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler.GetPVZs(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "ошибка получения")
}
