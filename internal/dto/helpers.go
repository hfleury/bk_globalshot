package dto

func ResponseSuccess(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
	}
}

func ResponseError(message string, errors []ErrorResponse) Response {
	return Response{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  errors,
	}
}

func ValidationError(field, message string, code ErrorCode) Response {
	return ResponseError("Validation failed.", []ErrorResponse{
		{
			Type:    "validation_error",
			Field:   field,
			Message: message,
			Code:    code,
		},
	})
}

func UnauthorizedResponse(message string) Response {
	if message == "" {
		message = ErrorCodeUnauthorized.DefaultMessage()
	}
	return ResponseError(message, []ErrorResponse{
		{
			Type:    "auth_error",
			Message: message,
			Code:    ErrorCodeUnauthorized,
		},
	})
}

func InternalServerErrorResponse() Response {
	return ResponseError("Internal server error.", []ErrorResponse{
		{
			Type:    "server_error",
			Message: ErrorCodeInternalServer.DefaultMessage(),
			Code:    ErrorCodeInternalServer,
		},
	})
}
