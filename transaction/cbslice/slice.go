package cbslice

import (
	"encoding/binary"
	"errors"
	"github.com/codingbeard/cbtransaction"
	"io"
	"sync"
)

var (
	DidNotReadEnoughDataTransactionSize = errors.New("did not read expected amount of data for the size of the transaction")
	DidNotReadEnoughData                = errors.New("did not read expected amount of data")
	NilSerialisedData                   = errors.New("nil serialised data")

	slicePool = &sync.Pool{
		New: func() interface{} {
			return New()
		},
	}
)

const (
	headerLength                = 25
	actionOffset                = 8
	transactionSizeByteLength   = 8
	encodingProviderKeyOffset   = 9
	encryptionProviderKeyOffset = 17
)

type Transaction []byte

func AcquireTransactionUnserialise(serialised []byte) *Transaction {
	transaction := slicePool.Get().(*Transaction)
	transaction.Unserialise(serialised)
	return transaction
}

func AcquireTransactionUnserialiseReader(reader io.Reader) (*Transaction, error) {
	transaction := slicePool.Get().(*Transaction)
	e := transaction.UnserialiseReader(reader)
	return transaction, e
}

func ReleaseTransaction(t *Transaction) {
	slicePool.Put(t)
}

func New() *Transaction {
	transaction := make(Transaction, headerLength)
	return &transaction
}

func NewFromReader(serialised io.Reader) (*Transaction, error) {
	transaction, e := NewUnserialiseReader(serialised)
	if e != nil {
		return nil, e
	}

	return transaction, nil
}

func (b *Transaction) SetTransactionId(transactionId uint64) {
	binary.LittleEndian.PutUint64(*b, transactionId)
}

func (b *Transaction) GetTransactionId() uint64 {
	return binary.LittleEndian.Uint64(*b)
}

func (b *Transaction) SetActionEnum(action cbtransaction.ActionEnum) {
	transaction := *b
	transaction[actionOffset] = byte(action)
	*b = transaction
}

func (b *Transaction) GetActionEnum() cbtransaction.ActionEnum {
	transaction := *b
	return cbtransaction.ActionEnum(transaction[actionOffset])
}

func (b *Transaction) SetEncodingProviderKey(key [8]byte) {
	transaction := *b
	prefix := append(
		transaction[:encodingProviderKeyOffset],
		key[0],
		key[1],
		key[2],
		key[3],
		key[4],
		key[5],
		key[6],
		key[7],
	)
	transaction = append(prefix, transaction[encryptionProviderKeyOffset:]...)
	*b = transaction
}

func (b *Transaction) GetEncodingProviderKey() [8]byte {
	transaction := *b
	return [8]byte{
		transaction[encodingProviderKeyOffset:encryptionProviderKeyOffset][0],
		transaction[encodingProviderKeyOffset:encryptionProviderKeyOffset][1],
		transaction[encodingProviderKeyOffset:encryptionProviderKeyOffset][2],
		transaction[encodingProviderKeyOffset:encryptionProviderKeyOffset][3],
		transaction[encodingProviderKeyOffset:encryptionProviderKeyOffset][4],
		transaction[encodingProviderKeyOffset:encryptionProviderKeyOffset][5],
		transaction[encodingProviderKeyOffset:encryptionProviderKeyOffset][6],
		transaction[encodingProviderKeyOffset:encryptionProviderKeyOffset][7],
	}
}

func (b *Transaction) SetEncryptionProviderKey(key [8]byte) {
	transaction := *b
	prefix := append(
		transaction[:encryptionProviderKeyOffset],
		key[0],
		key[1],
		key[2],
		key[3],
		key[4],
		key[5],
		key[6],
		key[7],
	)
	transaction = append(prefix, transaction[headerLength:]...)
	*b = transaction
}

func (b *Transaction) GetEncryptionProviderKey() [8]byte {
	transaction := *b
	return [8]byte{
		transaction[encryptionProviderKeyOffset:headerLength][0],
		transaction[encryptionProviderKeyOffset:headerLength][1],
		transaction[encryptionProviderKeyOffset:headerLength][2],
		transaction[encryptionProviderKeyOffset:headerLength][3],
		transaction[encryptionProviderKeyOffset:headerLength][4],
		transaction[encryptionProviderKeyOffset:headerLength][5],
		transaction[encryptionProviderKeyOffset:headerLength][6],
		transaction[encryptionProviderKeyOffset:headerLength][7],
	}
}

func (b *Transaction) SetData(data []byte) {
	*b = append(*b, data...)
}

func (b *Transaction) GetData() []byte {
	transaction := *b
	return transaction[headerLength:]
}

func (b *Transaction) GetLength() uint64 {
	return uint64(len(*b))
}

func (b *Transaction) Unserialise(transaction []byte) {
	*b = transaction
}

func (b *Transaction) UnserialiseReader(reader io.Reader) error {
	if reader != nil {
		transactionSizeBytes := make([]byte, transactionSizeByteLength)
		n, e := reader.Read(transactionSizeBytes)
		if e != nil {
			if errors.Is(e, io.EOF) {
				return DidNotReadEnoughDataTransactionSize
			} else {
				return e
			}
		}
		if n < transactionSizeByteLength {
			return DidNotReadEnoughDataTransactionSize
		} else {
			transactionSize := int(binary.LittleEndian.Uint64(transactionSizeBytes))
			if len(*b) != transactionSize {
				*b = make(Transaction, transactionSize)
			}
			n, e := reader.Read(*b)
			if e != nil {
				if errors.Is(e, io.EOF) && n != transactionSize {
					return DidNotReadEnoughData
				} else {
					return e
				}
			}
			if n < transactionSize {
				return DidNotReadEnoughData
			}
		}
	} else {
		return NilSerialisedData
	}

	return nil
}

func NewUnserialiseReader(reader io.Reader) (*Transaction, error) {
	if reader != nil {
		transactionSizeBytes := make([]byte, transactionSizeByteLength)
		n, e := reader.Read(transactionSizeBytes)
		if e != nil {
			if errors.Is(e, io.EOF) {
				return nil, DidNotReadEnoughDataTransactionSize
			} else {
				return nil, e
			}
		}
		if n < transactionSizeByteLength {
			return nil, DidNotReadEnoughDataTransactionSize
		} else {
			transactionSize := int(binary.LittleEndian.Uint64(transactionSizeBytes))

			transaction := make(Transaction, transactionSize)

			n, e := reader.Read(transaction)
			if e != nil {
				if errors.Is(e, io.EOF) && n != transactionSize {
					return nil, DidNotReadEnoughData
				} else {
					return nil, e
				}
			}
			if n < transactionSize {
				return nil, DidNotReadEnoughData
			}
			return &transaction, nil
		}
	} else {
		return nil, NilSerialisedData
	}
}

func (b *Transaction) Serialise() []byte {
	length := make([]byte, transactionSizeByteLength)
	binary.LittleEndian.PutUint64(length, b.GetLength())
	return append(length, *b...)
}

func (b *Transaction) SerialiseWriter(writer io.Writer) (n int, err error) {
	return writer.Write(*b)
}
