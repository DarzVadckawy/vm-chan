package domain

import "context"

// TextAnalysisRequest represents the input for text analysis
type TextAnalysisRequest struct {
	Sentence string `json:"sentence" binding:"required" example:"Hello world!"`
}

// TextAnalysisResponse represents the analysis result
type TextAnalysisResponse struct {
	Sentence       string `json:"sentence" example:"Hello world!"`
	WordCount      int    `json:"word_count" example:"2"`
	VowelCount     int    `json:"vowel_count" example:"3"`
	ConsonantCount int    `json:"consonant_count" example:"7"`
}

// User represents a user entity
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Never serialize password
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"password"`
}

// LoginResponse represents login response with JWT token
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in" example:"3600"`
}

// TextAnalysisService defines the business logic interface
type TextAnalysisService interface {
	AnalyzeText(ctx context.Context, sentence string) (*TextAnalysisResponse, error)
}

// AuthService defines authentication business logic interface
type AuthService interface {
	Login(ctx context.Context, username, password string) (*LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
}

// UserRepository defines user data access interface
type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
}
