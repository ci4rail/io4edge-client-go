/*
Copyright © 2024 Ci4Rail GmbH <engineering@ci4rail.com>
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

package errors

import (
	"fmt"
	"os"
)

// Er logs the error on stderr and terminates with exit code 1
func Er(msg interface{}) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", msg)
	os.Exit(1)
}

// ErrChk checks if err != nil. If so, prints the error and exists with exit code 1
func ErrChk(err error) {
	if err != nil {
		Er(err)
	}
}
