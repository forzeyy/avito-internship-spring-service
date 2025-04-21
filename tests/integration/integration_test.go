package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/forzeyy/avito-internship-spring-service/internal/database"
	"github.com/forzeyy/avito-internship-spring-service/internal/dto"
	"github.com/forzeyy/avito-internship-spring-service/internal/handlers"
	"github.com/forzeyy/avito-internship-spring-service/internal/middleware"
	"github.com/forzeyy/avito-internship-spring-service/internal/repos"
	"github.com/forzeyy/avito-internship-spring-service/internal/services"
	"github.com/forzeyy/avito-internship-spring-service/internal/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPVZAndReceptionIntegration(t *testing.T) {
	db, err := database.ConnectDatabase("postgres://postgres:postgres@localhost:5432/avito?sslmode=disable")
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	userRepo := repos.NewUserRepo(db)
	userSvc := services.NewUserService(userRepo, "secrettt", utils.DefaultAuthUtil{})
	authHandler := handlers.NewAuthHandler(userSvc)

	pvzRepo := repos.NewPVZRepo(db)
	pvzSvc := services.NewPVZService(pvzRepo)
	pvzHandler := handlers.NewPVZHandler(pvzSvc)

	receptionRepo := repos.NewReceptionRepo(db)
	receptionSvc := services.NewReceptionService(receptionRepo)
	receptionHandler := handlers.NewReceptionHandler(receptionSvc)

	productRepo := repos.NewProductRepo(db)
	productSvc := services.NewProductService(productRepo, receptionRepo)
	productHandler := handlers.NewProductHandler(productSvc)

	e := echo.New()
	RegisterHandlers(e, authHandler, pvzHandler, receptionHandler, productHandler)

	var moderatorToken string
	var employeeToken string
	var pvzID string

	t.Run("Dummy Login", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer([]byte(`{"role":"moderator"}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var tokenResp dto.Token
		err := json.Unmarshal(rec.Body.Bytes(), &tokenResp)
		assert.NoError(t, err)
		moderatorToken = tokenResp
		assert.NotEmpty(t, moderatorToken)

		req2 := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer([]byte(`{"role":"employee"}`)))
		req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec2 := httptest.NewRecorder()
		e.ServeHTTP(rec2, req2)

		assert.Equal(t, http.StatusOK, rec2.Code)

		err = json.Unmarshal(rec2.Body.Bytes(), &tokenResp)
		assert.NoError(t, err)
		employeeToken = tokenResp
		assert.NotEmpty(t, employeeToken)
	})

	t.Run("Create PVZ", func(t *testing.T) {
		body, _ := json.Marshal(handlers.CreatePVZRequest{
			City: dto.Москва,
		})
		req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+moderatorToken)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var createdPVZ dto.PVZ
		err := json.NewDecoder(rec.Body).Decode(&createdPVZ)
		assert.NoError(t, err)
		assert.NotNil(t, createdPVZ.Id)
		assert.Equal(t, "Москва", string(createdPVZ.City))

		pvzID = createdPVZ.Id.String()
	})

	t.Run("Create Reception", func(t *testing.T) {
		if pvzID == "" {
			t.Fatalf("ПВЗ не был создан")
		}

		reqBody := dto.PostReceptionsJSONRequestBody{
			PvzId: uuid.MustParse(pvzID),
		}

		reqBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(reqBytes))
		req.Header.Set("Authorization", "Bearer "+employeeToken)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var createdReception dto.Reception
		err := json.Unmarshal(rec.Body.Bytes(), &createdReception)
		assert.NoError(t, err)
		assert.Equal(t, "in_progress", string(createdReception.Status))
	})

	t.Run("Add 50 Products", func(t *testing.T) {
		if pvzID == "" {
			t.Fatalf("ПВЗ не был создан")
		}

		for i := 0; i < 50; i++ {
			reqBody := dto.PostProductsJSONRequestBody{
				Type:  "электроника",
				PvzId: uuid.MustParse(pvzID),
			}

			reqBytes, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(reqBytes))
			req.Header.Set("Authorization", "Bearer "+employeeToken)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("Close Reception", func(t *testing.T) {
		if pvzID == "" {
			t.Fatalf("ПВЗ не был создан")
		}

		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID+"/close_last_reception", nil)
		req.Header.Set("Authorization", "Bearer "+employeeToken)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var closedReception dto.Reception
		err := json.Unmarshal(rec.Body.Bytes(), &closedReception)
		assert.NoError(t, err)
		assert.Equal(t, "close", string(closedReception.Status))
	})
}

func RegisterHandlers(e *echo.Echo, authHandler *handlers.AuthHandler, pvzHandler *handlers.PVZHandler, receptionHandler *handlers.ReceptionHandler, productHandler *handlers.ProductHandler) {
	e.POST("/dummyLogin", authHandler.DummyLogin)
	e.POST("/register", authHandler.RegisterUser)
	e.POST("/login", authHandler.LoginUser)

	protected := e.Group("")
	protected.Use(middleware.JWTMiddleware("secrettt"), middleware.WithRole())

	protected.POST("/pvz", pvzHandler.CreatePVZ, middleware.OnlyModerator())
	protected.POST("/receptions", receptionHandler.CreateReception, middleware.OnlyEmployee())
	protected.POST("/pvz/:pvzId/close_last_reception", receptionHandler.CloseLastReception, middleware.OnlyEmployee())
	protected.POST("/pvz/:pvzId/delete_last_product", productHandler.DeleteLastProduct, middleware.OnlyEmployee())
	protected.POST("/products", productHandler.AddProduct, middleware.OnlyEmployee())
}
