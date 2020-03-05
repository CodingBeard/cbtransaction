package cbnone

import "io"

type Encryption struct {
}

func (n *Encryption) GetKey() [8]byte {
	return [8]byte{'n', 'o', 'n', 'e', 0, 0, 0, 0}
}

func (n *Encryption) Encrypt(data []byte) []byte {
	return data
}

func (n *Encryption) EncryptWriter(data []byte, writer io.Writer) error {
	_, e := writer.Write(data)
	return e
}

func (n *Encryption) Decrypt(encrypted []byte) []byte {
	return encrypted
}

func (n *Encryption) DecryptReader(reader io.Reader, out []byte) error {
	_, e := reader.Read(out)
	return e
}
