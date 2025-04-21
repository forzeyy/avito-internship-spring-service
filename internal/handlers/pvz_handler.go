package handlers

import (
	"net/http"

	"github.com/forzeyy/avito-internship-spring-service/internal/dto"
	"github.com/forzeyy/avito-internship-spring-service/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
)

type PVZHandler struct {
	pvzSvc services.PVZService
}

func NewPVZHandler(pvzSvc services.PVZService) *PVZHandler {
	return &PVZHandler{
		pvzSvc: pvzSvc,
	}
}

type CreatePVZRequest struct {
	City dto.PVZCity `json:"city"`
}

func (ph *PVZHandler) CreatePVZ(c echo.Context) error {
	var request CreatePVZRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: "невалидный запрос",
		})
	}

	pvz, err := ph.pvzSvc.CreatePVZ(c.Request().Context(), string(request.City))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.PVZ{
		Id:               (*types.UUID)(&pvz.ID),
		City:             dto.PVZCity(pvz.City),
		RegistrationDate: &pvz.RegDate,
	})
}

func (ph *PVZHandler) GetPVZs(c echo.Context) error {
	var params dto.GetPvzParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: "невалидный запрос",
		})
	}

	page := 1
	if params.Page != nil {
		page = *params.Page
	}

	limit := 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	pvzs, err := ph.pvzSvc.GetPVZs(c.Request().Context(), params.StartDate, params.EndDate, page, limit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: err.Error(),
		})
	}

	var dtoPvzs []dto.PVZ
	for _, pvz := range pvzs {
		dtoPvzs = append(dtoPvzs, dto.PVZ{
			Id:               (*types.UUID)(&pvz.ID),
			City:             dto.PVZCity(pvz.City),
			RegistrationDate: &pvz.RegDate,
		})
	}

	return c.JSON(http.StatusOK, dtoPvzs)
}
