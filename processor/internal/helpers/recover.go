package helpers

import (
	"errors"
	"fmt"
)

var ErrRecoveredFromPanic = errors.New("recovered from panic")

// RecoverRoutine recovers from a panic and sends the error to the provided error channel.
func RecoverRoutine(errChan chan<- error) {
	if r := recover(); r != nil {
		errChan <- fmt.Errorf("%w: %v", ErrRecoveredFromPanic, r)
	}
}
