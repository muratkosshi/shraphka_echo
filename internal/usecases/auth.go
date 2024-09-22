package usecases

import (
	"errors"
	"sharaphka_echo/internal/domain"
)

// AuthUseCase представляет сценарии использования для аутентификации
type AuthUseCase struct {
	userRepo   UserRepository
	jwtService JWTService
}

// UserRepository — интерфейс для работы с пользователями
type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
}

// JWTService — интерфейс для работы с JWT
type JWTService interface {
	GenerateToken(email string) (string, error)
	ValidateToken(tokenString string) (*domain.User, error)
}

// NewAuthUseCase создает новый AuthUseCase
func NewAuthUseCase(userRepo UserRepository, jwtService JWTService) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Login аутентифицирует пользователя и возвращает JWT токен
func (a *AuthUseCase) Login(email, password string) (string, error) {
	user, err := a.userRepo.FindByEmail(email)
	if err != nil || user == nil || user.Password != password {
		return "", errors.New("invalid credentials")
	}

	// Генерация JWT токена
	return a.jwtService.GenerateToken(user.Email)
}

// ValidateToken проверяет JWT токен и возвращает пользователя
func (a *AuthUseCase) ValidateToken(token string) (*domain.User, error) {
	return a.jwtService.ValidateToken(token)
}
