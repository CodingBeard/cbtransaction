package cbslice

import (
	"encoding/binary"
	"errors"
	"github.com/codingbeard/cbtransaction"
	"github.com/google/uuid"
	"io"
	"sync"
	"time"
)

var (
	Version1                            VersionEnum = 1
	DidNotReadEnoughDataTransactionSize             = errors.New("did not read expected amount of data for the size of the transaction")
	DidNotReadEnoughData                            = errors.New("did not read expected amount of data")
	NilSerialisedData                               = errors.New("nil serialised data")

	slicePool = &sync.Pool{
		New: func() interface{} {
			return NewVersion1()
		},
	}
)

const (
	transactionSizeByteLength = 8
	versionOffset             = 0

	headerLength1                = 34
	transactionIdOffset1         = 1
	actionOffset1                = 17
	encodingProviderKeyOffset1   = 18
	encryptionProviderKeyOffset1 = 26
)

type VersionEnum byte
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

func NewVersion1() *Transaction {
	transaction := make(Transaction, headerLength1)
	transaction[versionOffset] = byte(Version1)
	return &transaction
}

func NewFromReader(serialised io.Reader) (*Transaction, error) {
	transaction, e := NewUnserialiseReader(serialised)
	if e != nil {
		return nil, e
	}

	return transaction, nil
}

func (b *Transaction) SetVersion(version VersionEnum) {
	transaction := *b
	transaction[versionOffset] = byte(version)
	*b = transaction
}

func (b *Transaction) GetVersion() VersionEnum {
	transaction := *b
	return VersionEnum(transaction[versionOffset])
}

func (b *Transaction) SetTransactionId(transactionId uuid.UUID) {
	transaction := *b
	var offset int
	if b.GetVersion() == Version1 {
		offset = transactionIdOffset1
	}
	transaction[offset] = transactionId[0]
	transaction[offset+1] = transactionId[1]
	transaction[offset+2] = transactionId[2]
	transaction[offset+3] = transactionId[3]
	transaction[offset+4] = transactionId[4]
	transaction[offset+5] = transactionId[5]
	transaction[offset+6] = transactionId[6]
	transaction[offset+7] = transactionId[7]
	transaction[offset+8] = transactionId[8]
	transaction[offset+9] = transactionId[9]
	transaction[offset+10] = transactionId[10]
	transaction[offset+11] = transactionId[11]
	transaction[offset+12] = transactionId[12]
	transaction[offset+13] = transactionId[13]
	transaction[offset+14] = transactionId[14]
	transaction[offset+15] = transactionId[15]
	*b = transaction
}

func (b *Transaction) GetTransactionId() uuid.UUID {
	if b.GetVersion() == Version1 {
		transaction := *b

		transactionId := uuid.UUID{}
		transactionId[0] = transaction[transactionIdOffset1]
		transactionId[1] = transaction[transactionIdOffset1+1]
		transactionId[2] = transaction[transactionIdOffset1+2]
		transactionId[3] = transaction[transactionIdOffset1+3]
		transactionId[4] = transaction[transactionIdOffset1+4]
		transactionId[5] = transaction[transactionIdOffset1+5]
		transactionId[6] = transaction[transactionIdOffset1+6]
		transactionId[7] = transaction[transactionIdOffset1+7]
		transactionId[8] = transaction[transactionIdOffset1+8]
		transactionId[9] = transaction[transactionIdOffset1+9]
		transactionId[10] = transaction[transactionIdOffset1+10]
		transactionId[11] = transaction[transactionIdOffset1+11]
		transactionId[12] = transaction[transactionIdOffset1+12]
		transactionId[13] = transaction[transactionIdOffset1+13]
		transactionId[14] = transaction[transactionIdOffset1+14]
		transactionId[15] = transaction[transactionIdOffset1+15]

		return transactionId
	}
	return uuid.New()
}

func (b *Transaction) GetTime() time.Time {
	return time.Unix(b.GetTransactionId().Time().UnixTime())
}

func (b *Transaction) SetActionEnum(action cbtransaction.ActionEnum) {
	if b.GetVersion() == Version1 {
		transaction := *b
		transaction[actionOffset1] = byte(action)
		*b = transaction
	}
}

func (b *Transaction) GetActionEnum() cbtransaction.ActionEnum {
	if b.GetVersion() == Version1 {
		transaction := *b
		return cbtransaction.ActionEnum(transaction[actionOffset1])
	}

	return cbtransaction.ActionEnum(0)
}

func (b *Transaction) SetEncodingProviderKey(key [8]byte) {
	var encodingOffset, encryptionOffset int
	if b.GetVersion() == Version1 {
		encodingOffset = encodingProviderKeyOffset1
		encryptionOffset = encryptionProviderKeyOffset1
	}
	transaction := *b
	prefix := append(
		transaction[:encodingOffset],
		key[0],
		key[1],
		key[2],
		key[3],
		key[4],
		key[5],
		key[6],
		key[7],
	)
	transaction = append(prefix, transaction[encryptionOffset:]...)
	*b = transaction
}

func (b *Transaction) GetEncodingProviderKey() [8]byte {
	var encodingOffset, encryptionOffset int
	if b.GetVersion() == Version1 {
		encodingOffset = encodingProviderKeyOffset1
		encryptionOffset = encryptionProviderKeyOffset1
	}
	transaction := *b
	return [8]byte{
		transaction[encodingOffset:encryptionOffset][0],
		transaction[encodingOffset:encryptionOffset][1],
		transaction[encodingOffset:encryptionOffset][2],
		transaction[encodingOffset:encryptionOffset][3],
		transaction[encodingOffset:encryptionOffset][4],
		transaction[encodingOffset:encryptionOffset][5],
		transaction[encodingOffset:encryptionOffset][6],
		transaction[encodingOffset:encryptionOffset][7],
	}
}

func (b *Transaction) SetEncryptionProviderKey(key [8]byte) {
	var encryptionOffset, headerLength int
	if b.GetVersion() == Version1 {
		encryptionOffset = encryptionProviderKeyOffset1
		headerLength = headerLength1
	}
	transaction := *b
	prefix := append(
		transaction[:encryptionOffset],
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
	var encryptionOffset, headerLength int
	if b.GetVersion() == Version1 {
		encryptionOffset = encryptionProviderKeyOffset1
		headerLength = headerLength1
	}
	transaction := *b
	return [8]byte{
		transaction[encryptionOffset:headerLength][0],
		transaction[encryptionOffset:headerLength][1],
		transaction[encryptionOffset:headerLength][2],
		transaction[encryptionOffset:headerLength][3],
		transaction[encryptionOffset:headerLength][4],
		transaction[encryptionOffset:headerLength][5],
		transaction[encryptionOffset:headerLength][6],
		transaction[encryptionOffset:headerLength][7],
	}
}

func (b *Transaction) SetData(data []byte) {
	*b = append(*b, data...)
}

func (b *Transaction) GetData() []byte {
	var headerLength int
	if b.GetVersion() == Version1 {
		headerLength = headerLength1
	}
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
	length := make([]byte, transactionSizeByteLength, transactionSizeByteLength+len(*b))
	binary.LittleEndian.PutUint64(length, b.GetLength())
	return append(length, *b...)
}

func (b *Transaction) SerialiseWriter(writer io.Writer) (n int, err error) {
	return writer.Write(b.Serialise())
}
