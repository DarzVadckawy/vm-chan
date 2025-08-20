package handler

import (
	"net/http"

	"vm-chan/internal/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TextAnalysisHandler struct {
	service domain.TextAnalysisService
	logger  *zap.Logger
}

func NewTextAnalysisHandler(service domain.TextAnalysisService, logger *zap.Logger) *TextAnalysisHandler {
	return &TextAnalysisHandler{
		service: service,
		logger:  logger,
	}
}

func (h *TextAnalysisHandler) AnalyzeText(c *gin.Context) {
	var req domain.TextAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:       "Invalid request format",
			Code:        "validation_error",
			Description: "The request body does not match the expected format",
		})
		return
	}

	if req.Sentence == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:       "Sentence cannot be empty",
			Code:        "validation_error",
			Description: "The sentence field must not be empty",
		})
		return
	}

	result, err := h.service.AnalyzeText(c.Request.Context(), req.Sentence)
	if err != nil {
		h.logger.Error("Failed to analyze text", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error:       "Failed to analyze text",
			Code:        "internal_error",
			Description: "An unexpected error occurred while processing your request",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
