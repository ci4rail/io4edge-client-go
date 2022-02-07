/*
Copyright © 2022 Ci4Rail GmbH <engineering@ci4rail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

// HaveResponseStatus checks if err corresponds to functionblock status error code
func HaveResponseStatus(err error, status fbv1.Status) bool {
	if err == nil {
		return false
	}
	re, ok := err.(*ResponseError)
	if ok {
		if re.Status == status {
			return true
		}
	}
	return false
}
