package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"sharaphka_echo/internal/infrastructure/jwt"
	"sharaphka_echo/internal/usecases"
)

// RegisterRoutes регистрирует маршруты приложения
func RegisterRoutes(e *echo.Echo, authUseCase *usecases.AuthUseCase, jwtService *jwt.JWTService) {
	// Публичные маршруты
	e.POST("/login", LoginHandler(authUseCase))

	// JWT Middleware
	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(jwtService.Secret),
		ContextKey: "user",
		Claims:     &jwt.JWTClaim{},
	})

	// Защищенные маршруты
	protected := e.Group("/protected")
	protected.Use(jwtMiddleware)
	protected.GET("/profile", ProfileHandler(authUseCase))
}
