package domain

import "fmt"

type EntityValidationError struct {
	Field   string
	Message string
}

func (e *EntityValidationError) Error() string {
	return fmt.Sprintf("%s : %s", e.Field, e.Message)
}
