package commonerrors

import "fmt"

type InvalidEntityError struct {
	message string
}

func (_error *InvalidEntityError) Error() string {
	return fmt.Sprintln(_error.message)
}

func NewInvalidEntityError(message string) *InvalidEntityError {
	return &InvalidEntityError{message: message}
}
