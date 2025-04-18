package handlers

import (
	"net/http"

	"github.com/forzeyy/avito-internship-spring-service/internal/dto"
	"github.com/forzeyy/avito-internship-spring-service/internal/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
)

type ProductHandler struct {
	prodSvc services.ProductService
}

func NewProductHandler(prodSvc services.ProductService) *ProductHandler {
	return &ProductHandler{
		prodSvc: prodSvc,
	}
}

func (ph *ProductHandler) AddProduct(c echo.Context) error {
	var request dto.PostProductsJSONRequestBody
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: "невалидный запрос",
		})
	}

	if request.PvzId.String() == "" || request.Type == "" {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: "запрос должен содержать pvzId и type",
		})
	}

	product, err := ph.prodSvc.AddProduct(c.Request().Context(), string(request.Type), request.PvzId.String())
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, dto.Product{
		Id:          (*types.UUID)(&product.ID),
		Type:        dto.ProductType(product.Type),
		ReceptionId: (types.UUID)(product.ReceptionID),
		DateTime:    &product.DateTime,
	})
}

func (ph *ProductHandler) DeleteLastProduct(c echo.Context) error {
	pvzId := c.Param("pvzId")
	if _, err := uuid.Parse(pvzId); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: "неверный формат pvz_id",
		})
	}

	err := ph.prodSvc.DeleteLastProduct(c.Request().Context(), pvzId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error{
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}
