package bytestruct

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/codingbeard/cbtransaction"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"unsafe"
)

func TestNew(t *testing.T) {
	byteStruct := New()
	_ = byteStruct
}

func TestTransaction_SetTransactionId(t *testing.T) {
	byteStruct := New()
	transactionId := rand.Uint64()
	byteStruct.SetTransactionId(transactionId)
}

func TestTransaction_GetTransactionId(t *testing.T) {
	byteStruct := New()
	transactionId := rand.Uint64()
	byteStruct.SetTransactionId(transactionId)
	get := byteStruct.GetTransactionId()
	if get != transactionId {
		t.Errorf("transactionId was incorrect, got: %d, want: %d", get, transactionId)
	}
}

func TestTransaction_SetActionEnum(t *testing.T) {
	byteStruct := New()
	byteStruct.SetActionEnum(cbtransaction.ActionAdd)
	byteStruct.SetActionEnum(cbtransaction.ActionRemove)
	byteStruct.SetActionEnum(cbtransaction.ActionClear)
}

func TestTransaction_GetActionEnum(t *testing.T) {
	byteStruct := New()
	byteStruct.SetActionEnum(cbtransaction.ActionAdd)
	get := byteStruct.GetActionEnum()
	if get != cbtransaction.ActionAdd {
		t.Errorf("actionEnum was incorrect, got: %s, want: %s", string(get), string(cbtransaction.ActionAdd))
	}
	byteStruct.SetActionEnum(cbtransaction.ActionRemove)
	get = byteStruct.GetActionEnum()
	if get != cbtransaction.ActionRemove {
		t.Errorf("actionEnum was incorrect, got: %s, want: %s", string(get), string(cbtransaction.ActionRemove))
	}
	byteStruct.SetActionEnum(cbtransaction.ActionClear)
	get = byteStruct.GetActionEnum()
	if get != cbtransaction.ActionClear {
		t.Errorf("actionEnum was incorrect, got: %s, want: %s", string(get), string(cbtransaction.ActionClear))
	}
}

func TestTransaction_SetEncodingProviderKey(t *testing.T) {
	byteStruct := New()
	byteStruct.SetEncodingProviderKey([8]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'})
}

func TestTransaction_GetEncodingProviderKey(t *testing.T) {
	byteStruct := New()
	byteArray := [8]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
	byteStruct.SetEncodingProviderKey(byteArray)
	get := byteStruct.GetEncodingProviderKey()
	for key, arrayByte := range byteArray {
		if get[key] != arrayByte {
			t.Errorf("encodingProviderKey[%d] was incorrect, got: %s, want: %s", key, string(get[key]), string(arrayByte))
		}
	}
}

func TestTransaction_SetEncryptionProviderKey(t *testing.T) {
	byteStruct := New()
	byteStruct.SetEncryptionProviderKey([8]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'})
}

func TestTransaction_GetEncryptionProviderKey(t *testing.T) {
	byteStruct := New()
	byteArray := [8]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
	byteStruct.SetEncryptionProviderKey(byteArray)
	get := byteStruct.GetEncryptionProviderKey()
	for key, arrayByte := range byteArray {
		if get[key] != arrayByte {
			t.Errorf("encryptionProviderKey[%d] was incorrect, got: %s, want: %s", key, string(get[key]), string(arrayByte))
		}
	}
}

func TestTransaction_SetData(t *testing.T) {
	byteStruct := New()
	byteStruct.SetData([]byte("qwerty"))
}

func TestTransaction_GetData(t *testing.T) {
	byteStruct := New()
	data := []byte("qwerty")
	byteStruct.SetData(data)
	get := byteStruct.GetData()
	if bytes.Compare(data, get) != 0 {
		t.Errorf("data was incorrect, got: %b, want: %b", get, data)
	}
}

func TestTransaction_Serialise(t *testing.T) {
	byteStruct := New()
	t.Logf("new           %b %d", byteStruct, unsafe.Sizeof(*byteStruct))

	transactionId := uint64(0)
	byteStruct.SetTransactionId(transactionId)
	t.Logf("transactionId %b %d", byteStruct, unsafe.Sizeof(*byteStruct))

	byteStruct.SetActionEnum(cbtransaction.ActionAdd)
	t.Logf("action        %b %d", byteStruct, unsafe.Sizeof(*byteStruct))

	encodingProviderKey := [8]byte{0, 1, 2, 3, 4, 5, 6, 7}
	byteStruct.SetEncodingProviderKey(encodingProviderKey)
	t.Logf("encoding      %b %d", byteStruct, unsafe.Sizeof(*byteStruct))

	encryptionProviderKey := [8]byte{8, 9, 10, 11, 12, 13, 14, 15}
	byteStruct.SetEncryptionProviderKey(encryptionProviderKey)
	t.Logf("encryption    %b %d", byteStruct, unsafe.Sizeof(*byteStruct))

	data := []byte{16, 17, 18}
	byteStruct.SetData(data)
	t.Logf("data          %b %d", byteStruct, unsafe.Sizeof(*byteStruct))

	expected := []byte{
		28, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		0, 0, 0, 0, 0, 0, 0, 0, //transactionId
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
		16, 17, 18, //data
	}

	serialised := byteStruct.Serialise()
	if bytes.Compare(expected, serialised) != 0 {
		t.Errorf("serialised was incorrect, \ngot : %v, \nwant: %v", serialised, expected)
	}
}

func TestNewFromReader(t *testing.T) {
	reader := bytes.NewReader([]byte{
		28, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		0, 0, 0, 0, 0, 0, 0, 0, //transactionId
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
		16, 17, //data
	})
	byteStruct, e := NewFromReader(reader)
	if !errors.Is(e, DidNotReadEnoughData) {
		t.Errorf("e was incorrect, got: nil, want: %s", DidNotReadEnoughData.Error())
	}
	if byteStruct != nil {
		t.Errorf("byteStruct was incorrect, got: %s, want: nil", reflect.TypeOf(byteStruct).String())
	}

	reader = bytes.NewReader([]byte{
		28, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		0, 0, 0, 0, 0, 0, 0, 0, //transactionId
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
		16, 17, 18, //data
	})
	byteStruct, e = NewFromReader(reader)
	if e != nil {
		t.Errorf("e was incorrect, got: %s, want: nil", e.Error())
	}
	if byteStruct == nil {
		t.Errorf("byteStruct was incorrect, got: nil, want: Transaction")
		return
	}

	expectedTransactionId := uint64(0)
	getTransactionId := byteStruct.GetTransactionId()
	if getTransactionId != expectedTransactionId {
		t.Errorf("getTransactionId was incorrect, got: %d, want: %d", getTransactionId, expectedTransactionId)
	}

	expectedActionEnum := cbtransaction.ActionAdd
	getActionEnum := byteStruct.GetActionEnum()
	if getActionEnum != expectedActionEnum {
		t.Errorf("getActionEnum was incorrect, got: %s, want: %s", string(getActionEnum), string(expectedActionEnum))
	}

	expectedEncodingProviderKey := [8]byte{0, 1, 2, 3, 4, 5, 6, 7}
	getEncodingProviderKey := byteStruct.GetEncodingProviderKey()
	if getEncodingProviderKey != expectedEncodingProviderKey {
		t.Errorf("getEncodingProviderKey was incorrect, got: %v, want: %v", getEncodingProviderKey, expectedEncodingProviderKey)
	}

	expectedEncryptionProviderKey := [8]byte{8, 9, 10, 11, 12, 13, 14, 15}
	getEncryptionProviderKey := byteStruct.GetEncryptionProviderKey()
	if getEncryptionProviderKey != expectedEncryptionProviderKey {
		t.Errorf("getEncryptionProviderKey was incorrect, got: %v, want: %v", getEncryptionProviderKey, expectedEncryptionProviderKey)
	}

	expectedLength := uint64(28)
	getLength := byteStruct.GetLength()
	if getLength != expectedLength {
		t.Errorf("getLength was incorrect, got: %d, want: %d", getLength, expectedLength)
	}

	expectedData := []byte{16, 17, 18}
	getData := byteStruct.GetData()
	if bytes.Compare(getData, expectedData) != 0 {
		t.Errorf("getData was incorrect, got: %v, want: %v", getData, expectedData)
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

func BenchmarkTransaction_Serialise1KB(b *testing.B) {
	b.StopTimer()
	transactionBytes := []byte{
		0, 0, 0, 0, 0, 0, 0, 0, //transactionId
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
	}
	var data []byte
	for i := 0; i < 1024; i++ {
		data = append(data, 10)
	}
	size := make([]byte, 8)
	binary.LittleEndian.PutUint64(size, uint64(len(data)))

	transactionBytes = append(size, append(transactionBytes, data...)...)

	reader := bytes.NewReader(transactionBytes)
	byteStruct, e := NewFromReader(reader)
	if e != nil {
		b.Error(e)
		return
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = byteStruct.Serialise()
	}
}

//func BenchmarkReader1KBMemory(b *testing.B) {
//	b.StopTimer()
//	transactionBytes := []byte{
//		0, 0, 0, 0, 0, 0, 0, 0, //transactionId
//		43,                     //ActionAdd
//		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
//		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
//	}
//	var data []byte
//	for i := 0; i < 1024; i++ {
//		data = append(data, 10)
//	}
//	size := make([]byte, 8)
//	binary.LittleEndian.PutUint64(size, uint64(len(data)))
//
//	transactionBytes = append(transactionBytes, append(size, data...)...)
//
//	reader := bytes.NewReader(transactionBytes)
//
//	readTo1 := make([]byte, 33)
//	readTo2 := make([]byte, len(transactionBytes)-33)
//	b.StartTimer()
//	for i := 0; i < b.N; i++ {
//		b.StopTimer()
//		reader.Seek(0, 0)
//		b.StartTimer()
//		reader.Read(readTo1)
//		reader.Read(readTo2)
//	}
//}

func BenchmarkNewFromReader1KBMemory(b *testing.B) {
	b.StopTimer()
	transactionBytes := []byte{
		0, 0, 0, 0, 0, 0, 0, 0, //transactionId
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
	}
	var data []byte
	for i := 0; i < 1024; i++ {
		data = append(data, 10)
	}
	size := make([]byte, 8)
	binary.LittleEndian.PutUint64(size, uint64(len(data)))

	transactionBytes = append(size, append(transactionBytes, data...)...)

	reader := bytes.NewReader(transactionBytes)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		reader.Seek(0, 0)
		b.StartTimer()
		_, e := NewFromReader(reader)
		if e != nil {
			b.Error(e)
			return
		}
	}
}

func BenchmarkAcquireTransactionUnserialiseReader1KBMemory(b *testing.B) {
	b.StopTimer()
	transactionBytes := []byte{
		0, 0, 0, 0, 0, 0, 0, 0, //transactionId
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
	}
	var data []byte
	for i := 0; i < 1024; i++ {
		data = append(data, 10)
	}
	size := make([]byte, 8)
	binary.LittleEndian.PutUint64(size, uint64(len(data)))

	transactionBytes = append(size, append(transactionBytes, data...)...)

	reader := bytes.NewReader(transactionBytes)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		reader.Seek(0, 0)
		b.StartTimer()
		transaction, e := AcquireTransactionUnserialiseReader(reader)
		if e != nil {
			b.Error(e)
			return
		}
		ReleaseTransaction(transaction)
	}
}

func BenchmarkNewFromReader1KBFile(b *testing.B) {
	b.StopTimer()
	transactionBytes := []byte{
		0, 0, 0, 0, 0, 0, 0, 0, //transactionId
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
	}
	var data []byte
	for i := 0; i < 1024; i++ {
		data = append(data, 10)
	}
	size := make([]byte, 8)
	binary.LittleEndian.PutUint64(size, uint64(len(data)))

	transactionBytes = append(size, append(transactionBytes, data...)...)

	e := ioutil.WriteFile("BenchmarkNewFromReader1MBFile.data", transactionBytes, os.ModePerm)
	if e != nil {
		b.Error(e)
		return
	}

	file, e := os.OpenFile("BenchmarkNewFromReader1MBFile.data", os.O_RDONLY, os.ModePerm)
	if e != nil {
		b.Error(e)
		return
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		_, e := file.Seek(0, 0)
		if e != nil {
			b.Error(e)
			return
		}
		b.StartTimer()
		_, _ = NewFromReader(file)
	}
}

func BenchmarkAcquireTransactionUnserialiseReader1KBFile(b *testing.B) {
	b.StopTimer()
	transactionBytes := []byte{
		0, 0, 0, 0, 0, 0, 0, 0, //transactionId
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
	}
	var data []byte
	for i := 0; i < 1024; i++ {
		data = append(data, 10)
	}
	size := make([]byte, 8)
	binary.LittleEndian.PutUint64(size, uint64(len(data)))

	transactionBytes = append(size, append(transactionBytes, data...)...)

	e := ioutil.WriteFile("BenchmarkNewFromReader1MBFile.data", transactionBytes, os.ModePerm)
	if e != nil {
		b.Error(e)
		return
	}

	file, e := os.OpenFile("BenchmarkNewFromReader1MBFile.data", os.O_RDONLY, os.ModePerm)
	if e != nil {
		b.Error(e)
		return
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		_, e := file.Seek(0, 0)
		if e != nil {
			b.Error(e)
			return
		}
		b.StartTimer()
		transaction, e := AcquireTransactionUnserialiseReader(file)
		if e != nil {
			b.Error(e)
			return
		}
		ReleaseTransaction(transaction)
	}
}
