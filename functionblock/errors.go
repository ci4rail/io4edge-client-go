package functionblock

import (
	"fmt"

	fbv1 "github.com/ci4rail/io4edge_api/io4edge/go/functionblock/v1alpha1"
)

// ResponseError holds the details of functionblock errors
type ResponseError struct {
	Action        string
	Status        fbv1.Status
	FwErrorString string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("functionblock action '%s' failed with status %s: %s", e.Action, e.Status.String(), e.FwErrorString)
}

// StatusCode returns the numeric error status code from the functionblock
func (e *ResponseError) StatusCode() fbv1.Status {
	return e.Status
}

func responseErrorNew(action string, response *fbv1.Response) error {
	return &ResponseError{
		Status:        response.Status,
		FwErrorString: response.Error.String(),
		Action:        action,
	}
}
