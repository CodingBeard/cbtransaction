package cbtransaction

import (
	"io"
)

type ActionEnum byte

var (
	ActionAdd ActionEnum = '+'
	ActionRemove ActionEnum = '-'
	ActionClear ActionEnum = '*'
)

type Transaction interface {
	SetTransactionId(transactionId uint64)
	GetTransactionId() uint64
	SetActionEnum(action ActionEnum)
	GetActionEnum() ActionEnum
	SetEncodingProviderKey(key [8]byte)
	GetEncodingProviderKey() [8]byte
	SetEncryptionProviderKey(key [8]byte)
	GetEncryptionProviderKey() [8]byte
	SetData(data []byte)
	GetData() []byte
	GetLength() uint64
	Unserialise(transaction []byte)
	UnserialiseReader(reader io.Reader) error
	Serialise() []byte
	SerialiseWriter(writer io.Writer) (n int, err error)
}

func (e *ActionEnum) IsAdd() bool {
	return *e == ActionAdd
}

func (e *ActionEnum) IsRemove() bool {
	return *e == ActionRemove
}

func (e *ActionEnum) IsClear() bool {
	return *e == ActionClear
}
