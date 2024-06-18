package error

import "fmt"

type HTTPError struct {
	Description string `json:"description,omitempty"`
	StatusCode  int    `json:"statusCode"`
}

func (httpError HTTPError) Error() string {
	return fmt.Sprintf("description: %s", httpError.Description)
}

func NewHttpError(statusCode int) HTTPError {
	return HTTPError{
		StatusCode: statusCode,
	}
}

func NewHttpErrorWithDescription(statusCode int, description string) HTTPError {
	return HTTPError{
		StatusCode:  statusCode,
		Description: description,
	}
}
