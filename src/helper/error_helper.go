package helper

import (
	"errors"
	"fmt"
)

func NewError(msg string, args ...any) error {

	errorMsg := fmt.Sprintf(msg, args...)

	return errors.New(errorMsg)
}
