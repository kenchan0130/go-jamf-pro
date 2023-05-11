package utils

import (
	"fmt"
	"io"
	"log"
)

// HandleCloseFunc can be used to close an io.ReadCloser with message
func HandleCloseFunc(v io.ReadCloser, logger interface{}) {
	if err := v.Close(); err != nil {
		l := logger.(*log.Logger)
		if l != nil {
			l.Printf("Error closing io: %v", err)
			return
		}

		fmt.Printf("Error closing io: %v", err)
	}
}
