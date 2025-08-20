package repository

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"vm-chan/internal/domain"

	"go.uber.org/zap"
)

type userRepository struct {
	users  map[string]*domain.User
	logger *zap.Logger
}

func NewUserRepository(logger *zap.Logger) domain.UserRepository {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	users := map[string]*domain.User{
		"admin": {
			ID:       "1",
			Username: "admin",
			Password: string(hashedPassword),
		},
	}

	return &userRepository{
		users:  users,
		logger: logger,
	}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, exists := r.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}
