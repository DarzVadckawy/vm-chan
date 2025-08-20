package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTextAnalysisService_AnalyzeText(t *testing.T) {
	logger := zap.NewNop()
	service := NewTextAnalysisService(logger)

	tests := []struct {
		name               string
		sentence           string
		expectedWords      int
		expectedVowels     int
		expectedConsonants int
	}{
		{
			name:               "Simple sentence",
			sentence:           "Hello world",
			expectedWords:      2,
			expectedVowels:     3,
			expectedConsonants: 7,
		},
		{
			name:               "Single word",
			sentence:           "Hello",
			expectedWords:      1,
			expectedVowels:     2,
			expectedConsonants: 3,
		},
		{
			name:               "Empty sentence",
			sentence:           "",
			expectedWords:      0,
			expectedVowels:     0,
			expectedConsonants: 0,
		},
		{
			name:               "Sentence with punctuation",
			sentence:           "Hello, world!",
			expectedWords:      2,
			expectedVowels:     3,
			expectedConsonants: 7,
		},
		{
			name:               "Mixed case",
			sentence:           "HeLLo WoRLd",
			expectedWords:      2,
			expectedVowels:     3,
			expectedConsonants: 7,
		},
		{
			name:               "Numbers and special characters",
			sentence:           "Hello 123 world!",
			expectedWords:      3,
			expectedVowels:     3,
			expectedConsonants: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.AnalyzeText(context.Background(), tt.sentence)

			assert.NoError(t, err)
			assert.Equal(t, tt.sentence, result.Sentence)
			assert.Equal(t, tt.expectedWords, result.WordCount)
			assert.Equal(t, tt.expectedVowels, result.VowelCount)
			assert.Equal(t, tt.expectedConsonants, result.ConsonantCount)
		})
	}
}
