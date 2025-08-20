package service

import (
	"context"
	"testing"

	"vm-chan/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestAuthService_Login(t *testing.T) {
	logger := zap.NewNop()
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo, "test-secret", logger)

	t.Run("Successful login", func(t *testing.T) {
		// Use a test password hash instead of hardcoded production-like hash
		testPassword := "testpass123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
		user := &domain.User{
			ID:       "1",
			Username: "testuser",
			Password: string(hashedPassword),
		}

		mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(user, nil).Once()

		result, err := service.Login(context.Background(), "testuser", testPassword)

		assert.NoError(t, err)
		assert.NotEmpty(t, result.Token)
		assert.Greater(t, result.ExpiresIn, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid credentials", func(t *testing.T) {
		mockRepo.On("GetByUsername", mock.Anything, "nonexistent").Return((*domain.User)(nil), assert.AnError).Once()

		result, err := service.Login(context.Background(), "nonexistent", "password")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
