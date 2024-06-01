package helpers

type CustomError struct {
	Detail string
	Status int
}

func NotFoundError(detail string) *CustomError {
	return &CustomError{Detail: detail, Status: 404}
}

func AuthorizationError(detail string) *CustomError {
	return &CustomError{Detail: detail, Status: 403}
}

func AuthenticationError(detail string) *CustomError {
	return &CustomError{Detail: detail, Status: 401}
}

func ValidationError(detail string) *CustomError {
	return &CustomError{Detail: detail, Status: 400}
}
