package model

import (
	"regexp"

	"github.com/iris-contrib/parrot/parrot-api/errors"
)

// Validatable specifies the interface to validate structs.
type Validatable interface {
	Validate() error
}

var (
	emailRegex *regexp.Regexp
)
var (
	ErrValidationFailure = &errors.Error{
		Type:    "ValidationFailure",
		Message: "data validation failed"}
)

func init() {
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
}

// ValidEmail returns true if the string is of the valid email format.
func ValidEmail(str string) bool {
	return emailRegex.MatchString(str)
}

// HasMinLength returns true if the string's length is greater than or equal
// to the min parameter.
func HasMinLength(str string, min int) bool {
	return len(str) >= min
}

// NewValidationError constructs and returns a new error.
func NewValidationError(errs []errors.Error) error {
	return &errors.MultiError{
		Type:    ErrValidationFailure.Type,
		Message: ErrValidationFailure.Message,
		Errors:  errs}
}
