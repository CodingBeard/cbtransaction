package cbtransaction

import "io"

type Storage interface {
	Upload(filename string, reader io.Reader) error
	Download(filename string, writer io.Writer) error
	Concat(destination string, filenames ...string) error
	Delete(filename string) error
}
