package domain

import "context"

type TextAnalysisRequest struct {
	Sentence string `json:"sentence" binding:"required" example:"Hello world!"`
}

type TextAnalysisResponse struct {
	Sentence       string `json:"sentence" example:"Hello world!"`
	WordCount      int    `json:"word_count" example:"2"`
	VowelCount     int    `json:"vowel_count" example:"3"`
	ConsonantCount int    `json:"consonant_count" example:"7"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in" example:"3600"`
}

type TextAnalysisService interface {
	AnalyzeText(ctx context.Context, sentence string) (*TextAnalysisResponse, error)
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (*LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
}

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
}
