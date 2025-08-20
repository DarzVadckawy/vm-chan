package service

import (
	"context"
	"strings"
	"unicode"

	"vm-chan/internal/domain"

	"go.uber.org/zap"
)

type textAnalysisService struct {
	logger *zap.Logger
}

func NewTextAnalysisService(logger *zap.Logger) domain.TextAnalysisService {
	return &textAnalysisService{
		logger: logger,
	}
}

func (s *textAnalysisService) AnalyzeText(ctx context.Context, sentence string) (*domain.TextAnalysisResponse, error) {
	s.logger.Info("Analyzing text", zap.String("sentence", sentence))

	cleanSentence := strings.TrimSpace(sentence)
	if cleanSentence == "" {
		return &domain.TextAnalysisResponse{
			Sentence:       sentence,
			WordCount:      0,
			VowelCount:     0,
			ConsonantCount: 0,
		}, nil
	}

	words := strings.Fields(cleanSentence)
	wordCount := len(words)

	vowelCount := 0
	consonantCount := 0
	vowels := "aeiouAEIOU"

	for _, char := range cleanSentence {
		if unicode.IsLetter(char) {
			if strings.ContainsRune(vowels, char) {
				vowelCount++
			} else {
				consonantCount++
			}
		}
	}

	response := &domain.TextAnalysisResponse{
		Sentence:       sentence,
		WordCount:      wordCount,
		VowelCount:     vowelCount,
		ConsonantCount: consonantCount,
	}

	s.logger.Info("Text analysis completed",
		zap.String("sentence", sentence),
		zap.Int("words", wordCount),
		zap.Int("vowels", vowelCount),
		zap.Int("consonants", consonantCount),
	)

	return response, nil
}
