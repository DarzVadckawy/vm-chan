package domain

type ErrorResponse struct {
	Error       string `json:"error" example:"Invalid request format"`
	Code        string `json:"code,omitempty" example:"validation_error"`
	Description string `json:"description,omitempty" example:"The request body does not match the expected format"`
}
