package handlers

import (
	"net/http"

	"github.com/forzeyy/avito-internship-spring-service/internal/dto"
	"github.com/forzeyy/avito-internship-spring-service/internal/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
)

type ReceptionHandler struct {
	recSvc services.ReceptionService
}

func NewReceptionHandler(recSvc services.ReceptionService) *ReceptionHandler {
	return &ReceptionHandler{recSvc: recSvc}
}

func (rh *ReceptionHandler) CreateReception(c echo.Context) error {
	var request dto.PostReceptionsJSONRequestBody
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: "невалидный запрос",
		})
	}

	reception, err := rh.recSvc.CreateReception(c.Request().Context(), request.PvzId.String())
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, dto.Reception{
		Id:       (*types.UUID)(&reception.ID),
		PvzId:    (types.UUID)(reception.PVZID),
		Status:   dto.InProgress,
		DateTime: reception.DateTime,
	})
}

func (rh *ReceptionHandler) CloseLastReception(c echo.Context) error {
	pvzID := c.Param("pvzId")
	if _, err := uuid.Parse(pvzID); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: "неверный формат pvz_id",
		})
	}

	reception, err := rh.recSvc.CloseLastReception(c.Request().Context(), pvzID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.Reception{
		DateTime: reception.DateTime,
		Id:       (*types.UUID)(&reception.ID),
		PvzId:    (types.UUID)(reception.PVZID),
		Status:   dto.ReceptionStatus(reception.Status),
	})
}
