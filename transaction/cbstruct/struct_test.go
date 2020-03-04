package cbstruct

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
	tStruct := New()
	_ = tStruct
}

func TestTransaction_SetTransactionId(t *testing.T) {
	tStruct := New()
	transactionId := rand.Uint64()
	tStruct.SetTransactionId(transactionId)
}

func TestTransaction_GetTransactionId(t *testing.T) {
	tStruct := New()
	transactionId := rand.Uint64()
	tStruct.SetTransactionId(transactionId)
	get := tStruct.GetTransactionId()
	if get != transactionId {
		t.Errorf("transactionId was incorrect, got: %d, want: %d", get, transactionId)
	}
}

func TestTransaction_SetActionEnum(t *testing.T) {
	tStruct := New()
	tStruct.SetActionEnum(cbtransaction.ActionAdd)
	tStruct.SetActionEnum(cbtransaction.ActionRemove)
	tStruct.SetActionEnum(cbtransaction.ActionClear)
}

func TestTransaction_GetActionEnum(t *testing.T) {
	tStruct := New()
	tStruct.SetActionEnum(cbtransaction.ActionAdd)
	get := tStruct.GetActionEnum()
	if get != cbtransaction.ActionAdd {
		t.Errorf("actionEnum was incorrect, got: %s, want: %s", string(get), string(cbtransaction.ActionAdd))
	}
	tStruct.SetActionEnum(cbtransaction.ActionRemove)
	get = tStruct.GetActionEnum()
	if get != cbtransaction.ActionRemove {
		t.Errorf("actionEnum was incorrect, got: %s, want: %s", string(get), string(cbtransaction.ActionRemove))
	}
	tStruct.SetActionEnum(cbtransaction.ActionClear)
	get = tStruct.GetActionEnum()
	if get != cbtransaction.ActionClear {
		t.Errorf("actionEnum was incorrect, got: %s, want: %s", string(get), string(cbtransaction.ActionClear))
	}
}

func TestTransaction_SetEncodingProviderKey(t *testing.T) {
	tStruct := New()
	tStruct.SetEncodingProviderKey([8]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'})
}

func TestTransaction_GetEncodingProviderKey(t *testing.T) {
	tStruct := New()
	byteArray := [8]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
	tStruct.SetEncodingProviderKey(byteArray)
	get := tStruct.GetEncodingProviderKey()
	for key, arrayByte := range byteArray {
		if get[key] != arrayByte {
			t.Errorf("encodingProviderKey[%d] was incorrect, got: %s, want: %s", key, string(get[key]), string(arrayByte))
		}
	}
}

func TestTransaction_SetEncryptionProviderKey(t *testing.T) {
	tStruct := New()
	tStruct.SetEncryptionProviderKey([8]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'})
}

func TestTransaction_GetEncryptionProviderKey(t *testing.T) {
	tStruct := New()
	byteArray := [8]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
	tStruct.SetEncryptionProviderKey(byteArray)
	get := tStruct.GetEncryptionProviderKey()
	for key, arrayByte := range byteArray {
		if get[key] != arrayByte {
			t.Errorf("encryptionProviderKey[%d] was incorrect, got: %s, want: %s", key, string(get[key]), string(arrayByte))
		}
	}
}

func TestTransaction_SetData(t *testing.T) {
	tStruct := New()
	tStruct.SetData([]byte("qwerty"))
}

func TestTransaction_GetData(t *testing.T) {
	tStruct := New()
	data := []byte("qwerty")
	tStruct.SetData(data)
	get := tStruct.GetData()
	if bytes.Compare(data, get) != 0 {
		t.Errorf("data was incorrect, got: %b, want: %b", get, data)
	}
}

func TestTransaction_Serialise(t *testing.T) {
	tStruct := New()
	t.Logf("new           %b %d", tStruct, unsafe.Sizeof(*tStruct))

	transactionId := uint64(0)
	tStruct.SetTransactionId(transactionId)
	t.Logf("transactionId %b %d", tStruct, unsafe.Sizeof(*tStruct))

	tStruct.SetActionEnum(cbtransaction.ActionAdd)
	t.Logf("action        %b %d", tStruct, unsafe.Sizeof(*tStruct))

	encodingProviderKey := [8]byte{0, 1, 2, 3, 4, 5, 6, 7}
	tStruct.SetEncodingProviderKey(encodingProviderKey)
	t.Logf("encoding      %b %d", tStruct, unsafe.Sizeof(*tStruct))

	encryptionProviderKey := [8]byte{8, 9, 10, 11, 12, 13, 14, 15}
	tStruct.SetEncryptionProviderKey(encryptionProviderKey)
	t.Logf("encryption    %b %d", tStruct, unsafe.Sizeof(*tStruct))

	data := []byte{16, 17, 18}
	tStruct.SetData(data)
	t.Logf("data          %b %d", tStruct, unsafe.Sizeof(*tStruct))

	expected := []byte{
		28, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		0, 0, 0, 0, 0, 0, 0, 0, //transactionId
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
		16, 17, 18, //data
	}

	serialised := tStruct.Serialise()
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
	tStruct, e := NewFromReader(reader)
	if !errors.Is(e, DidNotReadEnoughData) {
		t.Errorf("e was incorrect, got: nil, want: %s", DidNotReadEnoughData.Error())
	}
	if tStruct != nil {
		t.Errorf("tStruct was incorrect, got: %s, want: nil", reflect.TypeOf(tStruct).String())
	}

	reader = bytes.NewReader([]byte{
		28, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		0, 0, 0, 0, 0, 0, 0, 0, //transactionId
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
		16, 17, 18, //data
	})
	tStruct, e = NewFromReader(reader)
	if e != nil {
		t.Errorf("e was incorrect, got: %s, want: nil", e.Error())
	}
	if tStruct == nil {
		t.Errorf("tStruct was incorrect, got: nil, want: Transaction")
		return
	}

	expectedTransactionId := uint64(0)
	getTransactionId := tStruct.GetTransactionId()
	if getTransactionId != expectedTransactionId {
		t.Errorf("getTransactionId was incorrect, got: %d, want: %d", getTransactionId, expectedTransactionId)
	}

	expectedActionEnum := cbtransaction.ActionAdd
	getActionEnum := tStruct.GetActionEnum()
	if getActionEnum != expectedActionEnum {
		t.Errorf("getActionEnum was incorrect, got: %s, want: %s", string(getActionEnum), string(expectedActionEnum))
	}

	expectedEncodingProviderKey := [8]byte{0, 1, 2, 3, 4, 5, 6, 7}
	getEncodingProviderKey := tStruct.GetEncodingProviderKey()
	if getEncodingProviderKey != expectedEncodingProviderKey {
		t.Errorf("getEncodingProviderKey was incorrect, got: %v, want: %v", getEncodingProviderKey, expectedEncodingProviderKey)
	}

	expectedEncryptionProviderKey := [8]byte{8, 9, 10, 11, 12, 13, 14, 15}
	getEncryptionProviderKey := tStruct.GetEncryptionProviderKey()
	if getEncryptionProviderKey != expectedEncryptionProviderKey {
		t.Errorf("getEncryptionProviderKey was incorrect, got: %v, want: %v", getEncryptionProviderKey, expectedEncryptionProviderKey)
	}

	expectedLength := uint64(28)
	getLength := tStruct.GetLength()
	if getLength != expectedLength {
		t.Errorf("getLength was incorrect, got: %d, want: %d", getLength, expectedLength)
	}

	expectedData := []byte{16, 17, 18}
	getData := tStruct.GetData()
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
	tStruct, e := NewFromReader(reader)
	if e != nil {
		b.Error(e)
		return
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = tStruct.Serialise()
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
