// Package fsrw implements application.FileReader by reading real source
// files.
package fsrw

import "os"

type Reader struct{}

func New() Reader { return Reader{} }

func (Reader) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}
