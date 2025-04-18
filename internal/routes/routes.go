package routes

import (
	"github.com/forzeyy/avito-internship-spring-service/internal/config"
	"github.com/forzeyy/avito-internship-spring-service/internal/database"
	"github.com/forzeyy/avito-internship-spring-service/internal/handlers"
	"github.com/forzeyy/avito-internship-spring-service/internal/middleware"
	"github.com/forzeyy/avito-internship-spring-service/internal/repos"
	"github.com/forzeyy/avito-internship-spring-service/internal/services"
	"github.com/forzeyy/avito-internship-spring-service/internal/utils"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, db *database.DB, cfg *config.Config) {
	// auth
	userRepo := repos.NewUserRepo(db)
	userSvc := services.NewUserService(userRepo, cfg.JWTSecret, utils.DefaultAuthUtil{})
	authHandler := handlers.NewAuthHandler(userSvc)

	// pvz
	pvzRepo := repos.NewPVZRepo(db)
	pvzSvc := services.NewPVZService(pvzRepo)
	pvzHandler := handlers.NewPVZHandler(pvzSvc)

	// reception
	receptionRepo := repos.NewReceptionRepo(db)
	receptionSvc := services.NewReceptionService(receptionRepo)
	receptionHandler := handlers.NewReceptionHandler(receptionSvc)

	// product
	productRepo := repos.NewProductRepo(db)
	productSvc := services.NewProductService(productRepo, receptionRepo)
	productHandler := handlers.NewProductHandler(productSvc)

	// open routes (auth)
	e.POST("/dummyLogin", authHandler.DummyLogin)
	e.POST("/login", authHandler.LoginUser)
	e.POST("/register", authHandler.RegisterUser)

	// protected routes
	protected := e.Group("")
	protected.Use(middleware.JWTMiddleware(cfg.JWTSecret), middleware.WithRole())

	// pvz
	protected.GET("/pvz", pvzHandler.GetPVZs)
	protected.POST("/pvz", pvzHandler.CreatePVZ, middleware.OnlyModerator())

	// reception
	protected.POST("/receptions", receptionHandler.CreateReception, middleware.OnlyEmployee())
	protected.POST("/pvz/:pvzId/close_last_reception", receptionHandler.CloseLastReception, middleware.OnlyEmployee())

	// product
	protected.POST("/products", productHandler.AddProduct, middleware.OnlyEmployee())
	protected.POST("/pvz/:pvzId/delete_last_product", productHandler.DeleteLastProduct, middleware.OnlyEmployee())
}
