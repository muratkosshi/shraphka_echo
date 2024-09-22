package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"sharaphka_echo/internal/domain"
	"time"
)

type JWTService struct {
	Secret string
}

// JWTClaim представляет пользовательские claims с валидацией
type JWTClaim struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Реализация метода Valid() для структуры JWTClaim
func (c *JWTClaim) Valid() error {
	if c.ExpiresAt == nil || !c.ExpiresAt.After(time.Now()) {
		return errors.New("token is expired")
	}
	return nil
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		Secret: secret,
	}
}

// GenerateToken создает JWT токен
func (j *JWTService) GenerateToken(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Создаем токен с использованием метода подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен с секретом
	return token.SignedString([]byte(j.Secret))
}

// ValidateToken проверяет валидность JWT токена и возвращает объект domain.User
func (j *JWTService) ValidateToken(tokenString string) (*domain.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		// Преобразование JWTClaim в domain.User
		user := &domain.User{
			Email: claims.Email,
		}
		return user, nil
	}

	return nil, errors.New("invalid token")
}
