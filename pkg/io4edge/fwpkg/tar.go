package fwpkg

import (
	"archive/tar"
	"errors"
	"fmt"
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
			fmt.Printf("reg %s\n", header.Name)
			if header.Name == fileName {
				if _, err := io.Copy(w, tr); err != nil {
					return err
				}
				return nil
			}
		}
	}
}
