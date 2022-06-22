package utils

import (
	"strings"

	"github.com/rl404/fairy/validation"
	"github.com/rl404/naka/internal/errors"
)

var val validation.Validator

func init() {
	val = validation.New(true)
	val.RegisterValidatorError("required", valErrRequired)
}

// Validate to validate struct using validate tag.
// Use pointer.
func Validate(data interface{}) error {
	return val.Validate(data)
}

func valErrRequired(f string, param ...string) error {
	return errors.ErrRequiredField(strings.ToLower(f))
}
