package repositories

import (
	"gorm.io/gorm"
	"log"
	"sharaphka_echo/internal/domain"
	"sharaphka_echo/internal/usecases"
)

type UserRepositoryDB struct {
	db *gorm.DB
}

func NewUserRepositoryDB(db *gorm.DB) usecases.UserRepository {
	return &UserRepositoryDB{db: db}
}

func (r *UserRepositoryDB) FindByEmail(email string) (*domain.User, error) {
	log.Printf("Executing query: SELECT * FROM users WHERE email = '%s'", email)

	var user domain.User
	// Используем GORM для поиска пользователя по email
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
