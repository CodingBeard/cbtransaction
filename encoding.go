package cbtransaction

import "io"

type Encoding interface {
	GetKey() [8]byte
	Encode(data interface{}) ([]byte, error)
	EncodeWriter(data interface{}, writer io.Writer) error
	Decode(encoded []byte) (interface{}, error)
	DecodeReader(reader io.Reader, out interface{}) error
}
