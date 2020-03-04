package byteslice

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

	byteSlicePool = &sync.Pool{
		New: func() interface{} {
			return New()
		},
	}
)

type Transaction []byte

func AcquireTransactionUnserialise(serialised []byte) *Transaction {
	transaction := byteSlicePool.Get().(*Transaction)
	transaction.Unserialise(serialised)
	return transaction
}

func AcquireTransactionUnserialiseReader(reader io.Reader) (*Transaction, error) {
	transaction := byteSlicePool.Get().(*Transaction)
	e := transaction.UnserialiseReader(reader)
	return transaction, e
}

func ReleaseTransaction(t *Transaction) {
	byteSlicePool.Put(t)
}

func New() *Transaction {
	transaction := make(Transaction, 25)
	return &transaction
}

func NewFromReader(serialised io.Reader) (*Transaction, error) {
	transaction := make(Transaction, 25)
	e := transaction.UnserialiseReader(serialised)
	if e != nil {
		return nil, e
	}

	return &transaction, nil
}

func (b *Transaction) SetTransactionId(transactionId uint64) {
	binary.LittleEndian.PutUint64(*b, transactionId)
}

func (b *Transaction) GetTransactionId() uint64 {
	return binary.LittleEndian.Uint64(*b)
}

func (b *Transaction) SetActionEnum(action cbtransaction.ActionEnum) {
	transaction := *b
	transaction[8] = byte(action)
	*b = transaction
}

func (b *Transaction) GetActionEnum() cbtransaction.ActionEnum {
	transaction := *b
	return cbtransaction.ActionEnum(transaction[8])
}

func (b *Transaction) SetEncodingProviderKey(key [8]byte) {
	transaction := *b
	prefix := append(
		transaction[:9],
		key[0],
		key[1],
		key[2],
		key[3],
		key[4],
		key[5],
		key[6],
		key[7],
	)
	transaction = append(prefix, transaction[17:]...)
	*b = transaction
}

func (b *Transaction) GetEncodingProviderKey() [8]byte {
	transaction := *b
	return [8]byte{
		transaction[9:17][0],
		transaction[9:17][1],
		transaction[9:17][2],
		transaction[9:17][3],
		transaction[9:17][4],
		transaction[9:17][5],
		transaction[9:17][6],
		transaction[9:17][7],
	}
}

func (b *Transaction) SetEncryptionProviderKey(key [8]byte) {
	transaction := *b
	prefix := append(
		transaction[:17],
		key[0],
		key[1],
		key[2],
		key[3],
		key[4],
		key[5],
		key[6],
		key[7],
	)
	transaction = append(prefix, transaction[25:]...)
	*b = transaction
}

func (b *Transaction) GetEncryptionProviderKey() [8]byte {
	transaction := *b
	return [8]byte{
		transaction[17:25][0],
		transaction[17:25][1],
		transaction[17:25][2],
		transaction[17:25][3],
		transaction[17:25][4],
		transaction[17:25][5],
		transaction[17:25][6],
		transaction[17:25][7],
	}
}

func (b *Transaction) SetData(data []byte) {
	*b = append(*b, data...)
}

func (b *Transaction) GetData() []byte {
	transaction := *b
	return transaction[25:]
}

func (b *Transaction) GetLength() uint64 {
	return uint64(len(*b))
}

func (b *Transaction) Unserialise(transaction []byte) {
	*b = transaction
}

func (b *Transaction) UnserialiseReader(reader io.Reader) error {
	if reader != nil {
		transactionSizeBytes := make([]byte, 8)
		n, e := reader.Read(transactionSizeBytes)
		if e != nil {
			if errors.Is(e, io.EOF) {
				return DidNotReadEnoughDataTransactionSize
			} else {
				return e
			}
		}
		if n < 8 {
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

func (b *Transaction) Serialise() []byte {
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, b.GetLength())
	return append(length, *b...)
}

func (b *Transaction) SerialiseWriter(writer io.Writer) (n int, err error) {
	return writer.Write(*b)
}
