package bytestruct

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

	byteStructPool = &sync.Pool{
		New: func() interface{} {
			return New()
		},
	}
)

type Transaction struct {
	transactionId         uint64
	actionEnum            cbtransaction.ActionEnum
	encodingProviderKey   [8]byte
	encryptionProviderKey [8]byte
	data                  []byte
}

func AcquireTransactionUnserialise(serialised []byte) *Transaction {
	transaction := byteStructPool.Get().(*Transaction)
	transaction.Unserialise(serialised)
	return transaction
}

func AcquireTransactionUnserialiseReader(reader io.Reader) (*Transaction, error) {
	transaction := byteStructPool.Get().(*Transaction)
	e := transaction.UnserialiseReader(reader)
	return transaction, e
}

func ReleaseTransaction(t *Transaction) {
	byteStructPool.Put(t)
}

func New() *Transaction {
	return &Transaction{}
}

func NewFromReader(serialised io.Reader) (*Transaction, error) {
	transaction := &Transaction{}
	e := transaction.UnserialiseReader(serialised)
	if e != nil {
		return nil, e
	}

	return transaction, nil
}

func (b *Transaction) SetTransactionId(transactionId uint64) {
	b.transactionId = transactionId
}

func (b *Transaction) GetTransactionId() uint64 {
	return b.transactionId
}

func (b *Transaction) SetActionEnum(action cbtransaction.ActionEnum) {
	b.actionEnum = action
}

func (b *Transaction) GetActionEnum() cbtransaction.ActionEnum {
	return b.actionEnum
}

func (b *Transaction) SetEncodingProviderKey(key [8]byte) {
	b.encodingProviderKey = key
}

func (b *Transaction) GetEncodingProviderKey() [8]byte {
	return b.encodingProviderKey
}

func (b *Transaction) SetEncryptionProviderKey(key [8]byte) {
	b.encryptionProviderKey = key
}

func (b *Transaction) GetEncryptionProviderKey() [8]byte {
	return b.encryptionProviderKey
}

func (b *Transaction) SetData(data []byte) {
	b.data = data
}

func (b *Transaction) GetData() []byte {
	return b.data
}

func (b *Transaction) GetLength() uint64 {
	return 8+1+8+8+uint64(len(b.data))
}

func (b *Transaction) Unserialise(transaction []byte) {
	b.transactionId = binary.LittleEndian.Uint64(transaction[:8])
	b.actionEnum = cbtransaction.ActionEnum(transaction[8])
	for i := 9; i < 17; i++ {
		b.encodingProviderKey[i-9] = transaction[i]
	}
	for i := 17; i < 25; i++ {
		b.encryptionProviderKey[i-17] = transaction[i]
	}
	b.data = transaction[25:]
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
			transactionBytes := make([]byte, transactionSize)
			n, e := reader.Read(transactionBytes)
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
			b.Unserialise(transactionBytes)
		}
	} else {
		return NilSerialisedData
	}

	return nil
}

func (b *Transaction) Serialise() []byte {
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, b.GetLength())

	transactionIdBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(transactionIdBytes, b.transactionId)
	encodingProviderKeyBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		encodingProviderKeyBytes[i] = b.encodingProviderKey[i]
	}
	encryptionProviderKeyBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		encryptionProviderKeyBytes[i] = b.encryptionProviderKey[i]
	}

	return append(length, append(transactionIdBytes, append([]byte{byte(b.actionEnum)}, append(encodingProviderKeyBytes, append(encryptionProviderKeyBytes, b.data...)...)...)...)...)
}

func (b *Transaction) SerialiseWriter(writer io.Writer) (n int, err error) {
	return writer.Write(b.Serialise())
}
