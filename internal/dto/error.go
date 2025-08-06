package dto

type ErrorResponse struct {
	Type    string    `json:"type"`
	Field   string    `json:"field"`
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

type ErrorCode string

const (
	ErrorCodeRequiredField     ErrorCode = "REQUIRED_FIELD"
	ErrorCodeInvalidFormat     ErrorCode = "INVALID_FORMAT"
	ErrorCodeUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden         ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound          ErrorCode = "NOT_FOUND"
	ErrorCodeValidationFailed  ErrorCode = "VALIDATION_FAILED"
	ErrorCodeInternalServer    ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeDuplicateEntry    ErrorCode = "DUPLICATE_ENTRY"
	ErrorCodeRateLimitExceeded ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrorCodeTimeout           ErrorCode = "TIMEOUT"
)

func (ec ErrorCode) DefaultMessage() string {
	switch ec {
	case ErrorCodeRequiredField:
		return "This field is required."
	case ErrorCodeInvalidFormat:
		return "The input format is invalid."
	case ErrorCodeUnauthorized:
		return "Authentication is required or has failed."
	case ErrorCodeForbidden:
		return "You don't have permission to perform this action."
	case ErrorCodeNotFound:
		return "The requested resource was not found."
	case ErrorCodeValidationFailed:
		return "One or more validation rules failed."
	case ErrorCodeInternalServer:
		return "An internal server error occurred."
	case ErrorCodeDuplicateEntry:
		return "This entry already exists."
	case ErrorCodeRateLimitExceeded:
		return "Too many requests. Please try again later."
	case ErrorCodeTimeout:
		return "The request timed out."
	default:
		return "An unexpected error occurred."
	}
}
