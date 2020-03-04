package cbbinary

import (
	"github.com/vmihailenco/msgpack"
	"io"
)

type Encoding struct {

}

type Config struct {

}

func New(config Config) *Encoding {
	return &Encoding{}
}

func (b *Encoding) GetKey() [8]byte {
	return [8]byte{'m', 's', 'g', 'p', 'a', 'c', 'k', 0}
}

func (b *Encoding) Encode(data interface{}) ([]byte, error) {
	return msgpack.Marshal(data)
}

func (b *Encoding) EncodeWriter(data interface{}, writer io.Writer) error {
	return msgpack.NewEncoder(writer).Encode(data)
}

func (b *Encoding) Decode(encoded []byte, out interface{}) error {
	return msgpack.Unmarshal(encoded, out)
}

func (b *Encoding) DecodeReader(reader io.Reader, out interface{}) error {
	return msgpack.NewDecoder(reader).Decode(out)
}



