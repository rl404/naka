package errors

import (
	"fmt"
)

// Error list.
var ()

// ErrRequiredField is error for missing field.
func ErrRequiredField(str string) error {
	return fmt.Errorf("required field %s", str)
}
