package error

import (
	"fmt"
)

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error Occured: %s\n", e.Message)
}

