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

// NewTextAnalysisService creates a new instance of text analysis service
func NewTextAnalysisService(logger *zap.Logger) domain.TextAnalysisService {
	return &textAnalysisService{
		logger: logger,
	}
}

// AnalyzeText analyzes the given sentence and returns word, vowel, and consonant counts
func (s *textAnalysisService) AnalyzeText(ctx context.Context, sentence string) (*domain.TextAnalysisResponse, error) {
	s.logger.Info("Analyzing text", zap.String("sentence", sentence))

	// Clean and process the sentence
	cleanSentence := strings.TrimSpace(sentence)
	if cleanSentence == "" {
		return &domain.TextAnalysisResponse{
			Sentence:       sentence,
			WordCount:      0,
			VowelCount:     0,
			ConsonantCount: 0,
		}, nil
	}

	// Count words
	words := strings.Fields(cleanSentence)
	wordCount := len(words)

	// Count vowels and consonants
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
