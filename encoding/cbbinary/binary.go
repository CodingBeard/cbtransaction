package cbbinary

import (
	"bytes"
	"encoding/binary"
	"io"
)

var KeyLittleEndian = [8]byte{'b', 'i', 'n', 'a', 'r', 'y', 'l', 0}
var KeyBigEndian = [8]byte{'b', 'i', 'n', 'a', 'r', 'y', 'b', 0}

type Encoding struct {
	endian binary.ByteOrder
}

type Config struct {
	Endian binary.ByteOrder
}

func New(config Config) *Encoding {
	return &Encoding{
		endian: config.Endian,
	}
}

func (b *Encoding) GetKey() [8]byte {
	var endianType byte
	if b.endian == binary.LittleEndian {
		return KeyLittleEndian
	} else if b.endian == binary.BigEndian {
		return KeyBigEndian
	}
	return [8]byte{'b', 'i', 'n', 'a', 'r', 'y', endianType, 0}
}

func (b *Encoding) Encode(data interface{}) ([]byte, error) {
	writer := bytes.NewBuffer([]byte{})
	e := binary.Write(writer, b.endian, data)
	if e != nil {
		return nil, e
	}
	return writer.Bytes(), nil
}

func (b *Encoding) EncodeWriter(data interface{}, writer io.Writer) error {
	return binary.Write(writer, b.endian, data)
}

func (b *Encoding) Decode(encoded []byte, out interface{}) error {
	reader := bytes.NewReader(encoded)
	return binary.Read(reader, b.endian, out)
}

func (b *Encoding) DecodeReader(reader io.Reader, out interface{}) error {
	return binary.Read(reader, b.endian, out)
}
