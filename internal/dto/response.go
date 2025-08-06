package dto

type Response struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    any             `json:"data,omitempty"`
	Errors  []ErrorResponse `json:"error,omitempty"`
}
