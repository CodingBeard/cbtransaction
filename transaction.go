package cbtransaction

import (
	"github.com/codingbeard/cbtransaction/transaction"
	"github.com/google/uuid"
	"io"
)

type Transaction interface {
	SetTransactionId(transactionId uuid.UUID)
	GetTransactionId() uuid.UUID
	SetActionEnum(action transaction.ActionEnum)
	GetActionEnum() transaction.ActionEnum
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
