package commands

import "fmt"

//InputError represents there having been a problem with the arguments given as
// the input of the command
type InputError struct {
	message string
}

func inputErrorf(format string, args ...interface{}) InputError {
	return InputError{message: fmt.Sprintf(format, args...)}
}

func (e InputError) Error() string {
	return e.message
}
