package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"sharaphka_echo/internal/infrastructure/jwt"
	"sharaphka_echo/internal/usecases"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler обрабатывает запросы на вход в систему
func LoginHandler(authUseCase *usecases.AuthUseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Вызов use case для аутентификации
		token, err := authUseCase.Login(req.Email, req.Password)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}

		// Возвращаем JWT токен
		return c.JSON(http.StatusOK, map[string]string{"token": token})
	}
}

// ProfileHandler обрабатывает запросы на получение профиля пользователя
func ProfileHandler(authUseCase *usecases.AuthUseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Получение данных пользователя из JWT токена
		userClaims := c.Get("user").(*jwt.JWTClaim)

		// Возвращаем данные профиля
		return c.JSON(http.StatusOK, map[string]string{
			"email":   userClaims.Email,
			"profile": "User profile data",
		})
	}
}
