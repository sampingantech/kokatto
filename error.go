package kokatto

import (
	"fmt"

	"github.com/pkg/errors"
)

// Error in kokatto client
var (
	ErrMissingParam = errors.New("missing params")
)

// Error from kokatto API
type Error struct {
	Status       string `json:"status"`
	StatusCode   string `json:"statusCode"`
	ErrorMessage string `json:"errorMessage"`
}

// Error return error message string
func (e *Error) Error() string {
	return fmt.Sprintf("code=%s status=%s message=%s", e.StatusCode, e.Status, e.ErrorMessage)
}
