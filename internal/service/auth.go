package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"

	"vm-chan/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type authService struct {
	userRepo  domain.UserRepository
	jwtSecret string
	logger    *zap.Logger
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new instance of authentication service
func NewAuthService(userRepo domain.UserRepository, jwtSecret string, logger *zap.Logger) domain.AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// Login authenticates user credentials and returns JWT token
func (s *authService) Login(ctx context.Context, username, password string) (*domain.LoginResponse, error) {
	s.logger.Info("Login attempt", zap.String("username", username))

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		s.logger.Error("User not found", zap.String("username", username), zap.Error(err))
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.Error("Invalid password", zap.String("username", username))
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "vm-chan",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Error(err))
		return nil, errors.New("failed to generate token")
	}

	s.logger.Info("Login successful", zap.String("username", username))

	return &domain.LoginResponse{
		Token:     tokenString,
		ExpiresIn: int(time.Until(expirationTime).Seconds()),
	}, nil
}

// ValidateToken validates JWT token and returns user information
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		s.logger.Error("Invalid token", zap.Error(err))
		return nil, errors.New("invalid token")
	}

	user := &domain.User{
		ID:       claims.UserID,
		Username: claims.Username,
	}

	return user, nil
}
