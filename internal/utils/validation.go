package utils

import (
	"fmt"
	"strings"

	"github.com/rl404/fairy/validation"
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
	return errRequiredField(strings.ToLower(f))
}

func errRequiredField(str string) error {
	return fmt.Errorf("required field %s", str)
}
