package main

import (
	"database/sql"
	"fmt"
	_ "go-crud-example/docs"
	"go-crud-example/internal/handler"
	"go-crud-example/internal/repository"
	"go-crud-example/internal/service"
	"go-crud-example/pkg/config"
	"go-crud-example/pkg/logger"
	"go-crud-example/pkg/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Users API
// @version 1.0
// @description API сервис для управления пользователями
// @host localhost:8000
// @BasePath /

func main() {

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем логгер
	logger := logger.NewLogger()

	// Инициализируем подключение к БД
	db, err := initDB(cfg.Database)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// Инициализируем слои приложения
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, logger)

	// Создаем роутер
	router := mux.NewRouter()

	// Добавляем middleware для логирования
	router.Use(middleware.LoggingMiddleware(logger))

	// Регистрируем маршруты
	userHandler.RegisterRoutes(router)

	// Добавляем Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Запускаем сервер
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Printf("Сервер запускается на порту %s...", cfg.Server.Port)
	logger.Fatal(http.ListenAndServe(serverAddr, router))
}

func initDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Создаем таблицу, если она не существует
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            age INT NOT NULL
        )
    `)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	return db, nil
}
