package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"sharaphka_echo/internal/infrastructure/database"
	"sharaphka_echo/internal/infrastructure/jwt"
	"sharaphka_echo/internal/interfaces/http"
	"sharaphka_echo/internal/interfaces/repositories"
	"sharaphka_echo/internal/usecases"
)

func main() {
	// Инициализация Echo
	e := echo.New()

	// Middleware для логов и восстановления после паники
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Инициализация подключения к базе данных
	db, err := database.NewDBConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Инициализация зависимостей
	userRepo := repositories.NewUserRepositoryDB(db)
	jwtService := jwt.NewJWTService("eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTcyNzAwNTY5OSwiaWF0IjoxNzI3MDA1Njk5fQ.DrV4D_lFleuay43oYqLl7aKneLGrnV02uFHWNLd2iRw")
	authUseCase := usecases.NewAuthUseCase(userRepo, jwtService)

	// Инициализация обработчиков HTTP
	http.RegisterRoutes(e, authUseCase, jwtService)

	// Запуск сервера
	log.Fatal(e.Start(":8080"))
}
