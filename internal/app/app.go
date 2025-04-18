package app

import (
	"fmt"

	"github.com/forzeyy/avito-internship-spring-service/internal/config"
	"github.com/forzeyy/avito-internship-spring-service/internal/database"
	"github.com/forzeyy/avito-internship-spring-service/internal/routes"
	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	dbConn, err := database.ConnectDatabase(dsn)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}
	defer dbConn.Close()

	e := echo.New()
	routes.InitRoutes(e, dbConn, cfg)
	e.Logger.Fatal(e.Start(":8080"))

	return nil
}
