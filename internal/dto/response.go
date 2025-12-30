package dto

type Response struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    any             `json:"data"`
	Errors  []ErrorResponse `json:"error,omitempty"`
}
