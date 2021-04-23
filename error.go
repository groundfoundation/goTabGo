package gotabgo

import (
	"fmt"
)

type ApiError struct {
	code    int
	message string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%d - %s", e.code, e.message)

}
