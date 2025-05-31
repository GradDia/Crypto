package dto

// ErrorResponseDto represents error response structure
// swagger:model ErrorResponseDto
type ErrorResponseDto struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
