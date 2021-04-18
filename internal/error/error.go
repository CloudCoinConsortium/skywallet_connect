package error

import (
	"fmt"
)

type Error struct {
	Code int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error Occured: %s\n", e.Message)
}

