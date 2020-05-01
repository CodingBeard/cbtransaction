package cbslice

import (
	"encoding/binary"
	"errors"
	"github.com/codingbeard/cbtransaction/transaction"
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
	tran := slicePool.Get().(*Transaction)
	tran.Unserialise(serialised)
	return tran
}

func AcquireTransactionUnserialiseReader(reader io.Reader) (*Transaction, error) {
	tran := slicePool.Get().(*Transaction)
	e := tran.UnserialiseReader(reader)
	return tran, e
}

func ReleaseTransaction(t *Transaction) {
	slicePool.Put(t)
}

func NewVersion1() *Transaction {
	tran := make(Transaction, headerLength1)
	tran[versionOffset] = byte(Version1)
	return &tran
}

func NewFromReader(serialised io.Reader) (*Transaction, error) {
	tran, e := NewUnserialiseReader(serialised)
	if e != nil {
		return nil, e
	}

	return tran, nil
}

func (b *Transaction) SetVersion(version VersionEnum) {
	tran := *b
	tran[versionOffset] = byte(version)
	*b = tran
}

func (b *Transaction) GetVersion() VersionEnum {
	tran := *b
	return VersionEnum(tran[versionOffset])
}

func (b *Transaction) SetTransactionId(transactionId uuid.UUID) {
	tran := *b
	var offset int
	if b.GetVersion() == Version1 {
		offset = transactionIdOffset1
	}
	tran[offset] = transactionId[0]
	tran[offset+1] = transactionId[1]
	tran[offset+2] = transactionId[2]
	tran[offset+3] = transactionId[3]
	tran[offset+4] = transactionId[4]
	tran[offset+5] = transactionId[5]
	tran[offset+6] = transactionId[6]
	tran[offset+7] = transactionId[7]
	tran[offset+8] = transactionId[8]
	tran[offset+9] = transactionId[9]
	tran[offset+10] = transactionId[10]
	tran[offset+11] = transactionId[11]
	tran[offset+12] = transactionId[12]
	tran[offset+13] = transactionId[13]
	tran[offset+14] = transactionId[14]
	tran[offset+15] = transactionId[15]
	*b = tran
}

func (b *Transaction) GetTransactionId() uuid.UUID {
	if b.GetVersion() == Version1 {
		tran := *b

		transactionId := uuid.UUID{}
		transactionId[0] = tran[transactionIdOffset1]
		transactionId[1] = tran[transactionIdOffset1+1]
		transactionId[2] = tran[transactionIdOffset1+2]
		transactionId[3] = tran[transactionIdOffset1+3]
		transactionId[4] = tran[transactionIdOffset1+4]
		transactionId[5] = tran[transactionIdOffset1+5]
		transactionId[6] = tran[transactionIdOffset1+6]
		transactionId[7] = tran[transactionIdOffset1+7]
		transactionId[8] = tran[transactionIdOffset1+8]
		transactionId[9] = tran[transactionIdOffset1+9]
		transactionId[10] = tran[transactionIdOffset1+10]
		transactionId[11] = tran[transactionIdOffset1+11]
		transactionId[12] = tran[transactionIdOffset1+12]
		transactionId[13] = tran[transactionIdOffset1+13]
		transactionId[14] = tran[transactionIdOffset1+14]
		transactionId[15] = tran[transactionIdOffset1+15]

		return transactionId
	}
	return uuid.New()
}

func (b *Transaction) GetTime() time.Time {
	return time.Unix(b.GetTransactionId().Time().UnixTime())
}

func (b *Transaction) SetActionEnum(action transaction.ActionEnum) {
	if b.GetVersion() == Version1 {
		tran := *b
		tran[actionOffset1] = byte(action)
		*b = tran
	}
}

func (b *Transaction) GetActionEnum() transaction.ActionEnum {
	if b.GetVersion() == Version1 {
		tran := *b
		return transaction.ActionEnum(tran[actionOffset1])
	}

	return transaction.ActionEnum(0)
}

func (b *Transaction) SetEncodingProviderKey(key [8]byte) {
	var encodingOffset, encryptionOffset int
	if b.GetVersion() == Version1 {
		encodingOffset = encodingProviderKeyOffset1
		encryptionOffset = encryptionProviderKeyOffset1
	}
	tran := *b
	prefix := append(
		tran[:encodingOffset],
		key[0],
		key[1],
		key[2],
		key[3],
		key[4],
		key[5],
		key[6],
		key[7],
	)
	tran = append(prefix, tran[encryptionOffset:]...)
	*b = tran
}

func (b *Transaction) GetEncodingProviderKey() [8]byte {
	var encodingOffset, encryptionOffset int
	if b.GetVersion() == Version1 {
		encodingOffset = encodingProviderKeyOffset1
		encryptionOffset = encryptionProviderKeyOffset1
	}
	tran := *b
	return [8]byte{
		tran[encodingOffset:encryptionOffset][0],
		tran[encodingOffset:encryptionOffset][1],
		tran[encodingOffset:encryptionOffset][2],
		tran[encodingOffset:encryptionOffset][3],
		tran[encodingOffset:encryptionOffset][4],
		tran[encodingOffset:encryptionOffset][5],
		tran[encodingOffset:encryptionOffset][6],
		tran[encodingOffset:encryptionOffset][7],
	}
}

func (b *Transaction) SetEncryptionProviderKey(key [8]byte) {
	var encryptionOffset, headerLength int
	if b.GetVersion() == Version1 {
		encryptionOffset = encryptionProviderKeyOffset1
		headerLength = headerLength1
	}
	tran := *b
	prefix := append(
		tran[:encryptionOffset],
		key[0],
		key[1],
		key[2],
		key[3],
		key[4],
		key[5],
		key[6],
		key[7],
	)
	tran = append(prefix, tran[headerLength:]...)
	*b = tran
}

func (b *Transaction) GetEncryptionProviderKey() [8]byte {
	var encryptionOffset, headerLength int
	if b.GetVersion() == Version1 {
		encryptionOffset = encryptionProviderKeyOffset1
		headerLength = headerLength1
	}
	tran := *b
	return [8]byte{
		tran[encryptionOffset:headerLength][0],
		tran[encryptionOffset:headerLength][1],
		tran[encryptionOffset:headerLength][2],
		tran[encryptionOffset:headerLength][3],
		tran[encryptionOffset:headerLength][4],
		tran[encryptionOffset:headerLength][5],
		tran[encryptionOffset:headerLength][6],
		tran[encryptionOffset:headerLength][7],
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
	tran := *b
	return tran[headerLength:]
}

func (b *Transaction) GetLength() uint64 {
	return uint64(len(*b))
}

func (b *Transaction) Unserialise(tran []byte) {
	*b = tran
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
