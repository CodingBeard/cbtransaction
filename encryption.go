package cbtransaction

import "io"

type Encryption interface {
	GetKey() [8]byte
	Encrypt(data []byte) []byte
	EncryptWriter(data []byte, writer io.Writer) error
	Decrypt(encrypted []byte) []byte
	DecryptReader(reader io.Reader, out []byte) error
}
