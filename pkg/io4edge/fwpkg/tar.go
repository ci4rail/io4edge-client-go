/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

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

package fwpkg

import (
	"archive/tar"
	"errors"
	"io"
)

// untarFileContent extracts the content of fileName from the uncompressed tar archive feed by r
// and writes the file content to w
// fileName must be the whole file name with path as in the archive e.g. "./test.txt"
func untarFileContent(r io.Reader, fileName string, w io.Writer) error {

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return errors.New("file " + fileName + " not found in archive")

		case err != nil:
			return err

		case header == nil:
			continue
		}

		if header.Typeflag == tar.TypeReg {
			if header.Name == fileName {
				if _, err := io.Copy(w, tr); err != nil {
					return err
				}
				return nil
			}
		}
	}
}
