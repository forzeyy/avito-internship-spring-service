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
	// Подключение к базе данных
	db, err := database.ConnectDatabase("postgres://postgres:postgres@localhost:5432/avito?sslmode=disable")
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Инициализация компонентов
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

	// Создаем экземпляр Echo для тестирования
	e := echo.New()

	// Регистрируем обработчики
	RegisterHandlers(e, authHandler, pvzHandler, receptionHandler, productHandler)

	// Переменная для хранения ID созданного ПВЗ
	var pvzID string

	// Шаг 1: Создание нового ПВЗ
	t.Run("Create PVZ", func(t *testing.T) {
		reqBody := dto.PostPvzJSONRequestBody{
			City: dto.PVZCity("Москва"),
		}

		reqBodyBytes, _ := json.Marshal(reqBody)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Authorization", "Bearer moderatorToken") // Модераторский токен
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var createdPVZ dto.PVZ
		err := json.Unmarshal(rec.Body.Bytes(), &createdPVZ)
		assert.NoError(t, err)
		assert.NotEmpty(t, createdPVZ.Id)
		assert.Equal(t, "Москва", string(createdPVZ.City))

		// Сохраняем ID созданного ПВЗ для дальнейшего использования
		pvzID = createdPVZ.Id.String()
		t.Logf("Created PVZ with ID: %s", pvzID)
	})

	// Шаг 2: Добавление новой приемки заказов
	t.Run("Create Reception", func(t *testing.T) {
		if pvzID == "" {
			t.Fatalf("ПВЗ не был создан")
		}

		reqBody := dto.PostReceptionsJSONRequestBody{
			PvzId: uuid.MustParse(pvzID),
		}

		reqBodyBytes, _ := json.Marshal(reqBody)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Authorization", "Bearer employeeToken")
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var createdReception dto.Reception
		err := json.Unmarshal(rec.Body.Bytes(), &createdReception)
		assert.NoError(t, err)
		assert.Equal(t, "in_progress", string(createdReception.Status))
		t.Logf("Created Reception with ID: %s", createdReception.Id)
	})

	// Шаг 3: Добавление 50 товаров в текущую приемку
	t.Run("Add 50 Products", func(t *testing.T) {
		if pvzID == "" {
			t.Fatalf("ПВЗ не был создан")
		}

		for i := 0; i < 50; i++ {
			reqBody := dto.PostProductsJSONRequestBody{
				Type:  "электроника",
				PvzId: uuid.MustParse(pvzID),
			}

			reqBodyBytes, _ := json.Marshal(reqBody)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(reqBodyBytes))
			req.Header.Set("Authorization", "Bearer employeeToken")
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusCreated, rec.Code)
		}

		t.Logf("Added 50 products to the reception")
	})

	// Шаг 4: Закрытие приемки заказов
	t.Run("Close Reception", func(t *testing.T) {
		if pvzID == "" {
			t.Fatalf("ПВЗ не был создан")
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID+"/close_last_reception", nil)
		req.Header.Set("Authorization", "Bearer employeeToken")

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var closedReception dto.Reception
		err := json.Unmarshal(rec.Body.Bytes(), &closedReception)
		assert.NoError(t, err)
		assert.Equal(t, "close", string(closedReception.Status))
		t.Logf("Closed reception with ID: %s", closedReception.Id)
	})
}

// Функция для регистрации обработчиков
func RegisterHandlers(e *echo.Echo, authHandler *handlers.AuthHandler, pvzHandler *handlers.PVZHandler, receptionHandler *handlers.ReceptionHandler, productHandler *handlers.ProductHandler) {
	// Авторизация
	e.POST("/dummyLogin", authHandler.DummyLogin)
	e.POST("/register", authHandler.RegisterUser)
	e.POST("/login", authHandler.LoginUser)

	// Группа защищенных маршрутов
	protected := e.Group("")
	protected.Use(middleware.JWTMiddleware("secrettt"), middleware.WithRole())

	// ПВЗ (создание доступно только модераторам)
	protected.POST("/pvz", pvzHandler.CreatePVZ, middleware.OnlyModerator())

	// Приемка товаров (доступно только сотрудникам ПВЗ)
	protected.POST("/receptions", receptionHandler.CreateReception, middleware.OnlyEmployee())
	protected.POST("/pvz/:pvzId/close_last_reception", receptionHandler.CloseLastReception, middleware.OnlyEmployee())

	// Товары (доступно только сотрудникам ПВЗ)
	protected.POST("/pvz/:pvzId/delete_last_product", productHandler.DeleteLastProduct, middleware.OnlyEmployee())
	protected.POST("/products", productHandler.AddProduct, middleware.OnlyEmployee())
}
