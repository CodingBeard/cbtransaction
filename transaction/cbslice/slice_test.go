package cbslice

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/codingbeard/cbtransaction"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	slice := NewVersion1()
	_ = slice
}

func TestTransaction_SetVersion(t *testing.T) {
	slice := NewVersion1()
	version := Version1
	slice.SetVersion(version)
}

func TestTransaction_GetVersion(t *testing.T) {
	slice := NewVersion1()
	version := Version1
	slice.SetVersion(version)
	get := slice.GetVersion()
	if get != version {
		t.Errorf("version was incorrect, got: %d, want: %d", get, version)
	}
}

func TestTransaction_SetTransactionId(t *testing.T) {
	slice := NewVersion1()

	transactionId, e := uuid.NewUUID()
	if e != nil {
		t.Error(e)
	}
	slice.SetTransactionId(transactionId)
}

func TestTransaction_GetTransactionId(t *testing.T) {
	slice := NewVersion1()

	transactionId, e := uuid.NewUUID()
	if e != nil {
		t.Error(e)
	}
	slice.SetTransactionId(transactionId)
	get := slice.GetTransactionId()
	if get.String() != transactionId.String() {
		t.Errorf("transactionId was incorrect, got: %s, want: %s", get.String(), transactionId.String())
	}
}

func TestTransaction_GetTime(t *testing.T) {
	slice := NewVersion1()

	transactionId, e := uuid.NewUUID()
	if e != nil {
		t.Error(e)
	}
	slice.SetTransactionId(transactionId)
	get := slice.GetTime()
	unix, _ := transactionId.Time().UnixTime()
	if get.Unix() != unix {
		t.Errorf("transactionId was incorrect, got: %d, want: %d", get.Unix(), unix)
	}
}

func TestTransaction_SetActionEnum(t *testing.T) {
	slice := NewVersion1()

	slice.SetActionEnum(cbtransaction.ActionAdd)
	slice.SetActionEnum(cbtransaction.ActionRemove)
	slice.SetActionEnum(cbtransaction.ActionClear)
}

func TestTransaction_GetActionEnum(t *testing.T) {
	slice := NewVersion1()

	slice.SetActionEnum(cbtransaction.ActionAdd)
	get := slice.GetActionEnum()
	if get != cbtransaction.ActionAdd {
		t.Errorf("actionEnum was incorrect, got: %s, want: %s", string(get), string(cbtransaction.ActionAdd))
	}
	slice.SetActionEnum(cbtransaction.ActionRemove)
	get = slice.GetActionEnum()
	if get != cbtransaction.ActionRemove {
		t.Errorf("actionEnum was incorrect, got: %s, want: %s", string(get), string(cbtransaction.ActionRemove))
	}
	slice.SetActionEnum(cbtransaction.ActionClear)
	get = slice.GetActionEnum()
	if get != cbtransaction.ActionClear {
		t.Errorf("actionEnum was incorrect, got: %s, want: %s", string(get), string(cbtransaction.ActionClear))
	}
}

func TestTransaction_SetEncodingProviderKey(t *testing.T) {
	slice := NewVersion1()

	slice.SetEncodingProviderKey([8]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'})
}

func TestTransaction_GetEncodingProviderKey(t *testing.T) {
	slice := NewVersion1()

	byteArray := [8]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
	slice.SetEncodingProviderKey(byteArray)
	get := slice.GetEncodingProviderKey()
	for key, arrayByte := range byteArray {
		if get[key] != arrayByte {
			t.Errorf("encodingProviderKey[%d] was incorrect, got: %s, want: %s", key, string(get[key]), string(arrayByte))
		}
	}
}

func TestTransaction_SetEncryptionProviderKey(t *testing.T) {
	slice := NewVersion1()

	slice.SetEncryptionProviderKey([8]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'})
}

func TestTransaction_GetEncryptionProviderKey(t *testing.T) {
	slice := NewVersion1()

	byteArray := [8]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
	slice.SetEncryptionProviderKey(byteArray)
	get := slice.GetEncryptionProviderKey()
	for key, arrayByte := range byteArray {
		if get[key] != arrayByte {
			t.Errorf("encryptionProviderKey[%d] was incorrect, got: %s, want: %s", key, string(get[key]), string(arrayByte))
		}
	}
}

func TestTransaction_SetData(t *testing.T) {
	slice := NewVersion1()

	slice.SetData([]byte("qwerty"))
}

func TestTransaction_GetData(t *testing.T) {
	slice := NewVersion1()

	data := []byte("qwerty")
	slice.SetData(data)
	get := slice.GetData()
	if bytes.Compare(data, get) != 0 {
		t.Errorf("data was incorrect, got: %b, want: %b", get, data)
	}
}

func TestTransaction_Serialise(t *testing.T) {
	slice := NewVersion1()

	t.Logf("new           %b %d", slice, len(*slice))

	transactionId := uuid.New()
	slice.SetTransactionId(transactionId)
	t.Logf("transactionId %b %d", slice, len(*slice))

	slice.SetActionEnum(cbtransaction.ActionAdd)
	t.Logf("action        %b %d", slice, len(*slice))

	encodingProviderKey := [8]byte{0, 1, 2, 3, 4, 5, 6, 7}
	slice.SetEncodingProviderKey(encodingProviderKey)
	t.Logf("encoding      %b %d", slice, len(*slice))

	encryptionProviderKey := [8]byte{8, 9, 10, 11, 12, 13, 14, 15}
	slice.SetEncryptionProviderKey(encryptionProviderKey)
	t.Logf("encryption    %b %d", slice, len(*slice))

	data := []byte{16, 17, 18}
	slice.SetData(data)
	t.Logf("data          %b %d", slice, len(*slice))

	expected := []byte{
		37, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		byte(Version1),                                                         //Version1
		transactionId[0], transactionId[1], transactionId[2], transactionId[3], //UUID
		transactionId[4], transactionId[5], transactionId[6], transactionId[7], //UUID
		transactionId[8], transactionId[9], transactionId[10], transactionId[11], //UUID
		transactionId[12], transactionId[13], transactionId[14], transactionId[15], //UUID
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
		16, 17, 18, //data
	}

	serialised := slice.Serialise()
	if bytes.Compare(expected, serialised) != 0 {
		t.Errorf("serialised was incorrect, \ngot : %v, \nwant: %v", serialised, expected)
	}
}

func TestTransaction_SerialiseWriter(t *testing.T) {
	slice := NewVersion1()

	t.Logf("new           %b %d", slice, len(*slice))

	transactionId := uuid.New()
	slice.SetTransactionId(transactionId)
	t.Logf("transactionId %b %d", slice, len(*slice))

	slice.SetActionEnum(cbtransaction.ActionAdd)
	t.Logf("action        %b %d", slice, len(*slice))

	encodingProviderKey := [8]byte{0, 1, 2, 3, 4, 5, 6, 7}
	slice.SetEncodingProviderKey(encodingProviderKey)
	t.Logf("encoding      %b %d", slice, len(*slice))

	encryptionProviderKey := [8]byte{8, 9, 10, 11, 12, 13, 14, 15}
	slice.SetEncryptionProviderKey(encryptionProviderKey)
	t.Logf("encryption    %b %d", slice, len(*slice))

	data := []byte{16, 17, 18}
	slice.SetData(data)
	t.Logf("data          %b %d", slice, len(*slice))

	expected := []byte{
		37, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		byte(Version1),                                                         //Version1
		transactionId[0], transactionId[1], transactionId[2], transactionId[3], //UUID
		transactionId[4], transactionId[5], transactionId[6], transactionId[7], //UUID
		transactionId[8], transactionId[9], transactionId[10], transactionId[11], //UUID
		transactionId[12], transactionId[13], transactionId[14], transactionId[15], //UUID
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
		16, 17, 18, //data
	}

	writer := bytes.NewBuffer([]byte{})

	n, e := slice.SerialiseWriter(writer)
	if e != nil {
		t.Errorf("e was incorrect, got : %v, want: nil", e)
	}
	if n != 37+8 {
		t.Errorf("n was incorrect, got : %d, want: 37", n)
	}
	if bytes.Compare(expected, writer.Bytes()) != 0 {
		t.Errorf("serialised was incorrect, \ngot : %v, \nwant: %v", writer.Bytes(), expected)
	}
}

func TestNewFromReader(t *testing.T) {
	reader := bytes.NewReader([]byte{
		37, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		byte(Version1),                                 //Version1
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, //UUID
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
		16, 17, //data
	})
	slice, e := NewFromReader(reader)
	if !errors.Is(e, DidNotReadEnoughData) {
		t.Errorf("e was incorrect, got: nil, want: %s", DidNotReadEnoughData.Error())
	}
	if slice != nil {
		t.Errorf("slice was incorrect, got: %s, want: nil", reflect.TypeOf(slice).String())
	}

	reader = bytes.NewReader([]byte{
		37, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		byte(Version1),                                 //Version1
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, //UUID
		43,                     //ActionAdd
		0, 1, 2, 3, 4, 5, 6, 7, //encodingProviderKey
		8, 9, 10, 11, 12, 13, 14, 15, //encryptionProviderKey
		16, 17, 18, //data
	})
	slice, e = NewFromReader(reader)
	if e != nil {
		t.Errorf("e was incorrect, got: %s, want: nil", e.Error())
	}
	if slice == nil {
		t.Errorf("slice was incorrect, got: nil, want: Transaction")
		return
	}

	expectedTransactionId := uuid.UUID{}
	getTransactionId := slice.GetTransactionId()
	if getTransactionId.String() != expectedTransactionId.String() {
		t.Errorf("getTransactionId was incorrect, got: %s, want: %s", getTransactionId.String(), expectedTransactionId.String())
	}

	expectedActionEnum := cbtransaction.ActionAdd
	getActionEnum := slice.GetActionEnum()
	if getActionEnum != expectedActionEnum {
		t.Errorf("getActionEnum was incorrect, got: %s, want: %s", string(getActionEnum), string(expectedActionEnum))
	}

	expectedEncodingProviderKey := [8]byte{0, 1, 2, 3, 4, 5, 6, 7}
	getEncodingProviderKey := slice.GetEncodingProviderKey()
	if getEncodingProviderKey != expectedEncodingProviderKey {
		t.Errorf("getEncodingProviderKey was incorrect, got: %v, want: %v", getEncodingProviderKey, expectedEncodingProviderKey)
	}

	expectedEncryptionProviderKey := [8]byte{8, 9, 10, 11, 12, 13, 14, 15}
	getEncryptionProviderKey := slice.GetEncryptionProviderKey()
	if getEncryptionProviderKey != expectedEncryptionProviderKey {
		t.Errorf("getEncryptionProviderKey was incorrect, got: %v, want: %v", getEncryptionProviderKey, expectedEncryptionProviderKey)
	}

	expectedLength := uint64(37)
	getLength := slice.GetLength()
	if getLength != expectedLength {
		t.Errorf("getLength was incorrect, got: %d, want: %d", getLength, expectedLength)
	}

	expectedData := []byte{16, 17, 18}
	getData := slice.GetData()
	if bytes.Compare(getData, expectedData) != 0 {
		t.Errorf("getData was incorrect, got: %v, want: %v", getData, expectedData)
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewVersion1()
	}
}

func BenchmarkTransaction_Serialise1KB(b *testing.B) {
	b.StopTimer()
	transactionBytes := []byte{
		37, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		byte(Version1),                                 //Version1
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, //UUID
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
	slice, e := NewFromReader(reader)
	if e != nil {
		b.Error(e)
		return
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = slice.Serialise()
	}
}

func BenchmarkTransaction_SerialiseWriter1KB(b *testing.B) {
	b.StopTimer()
	transactionBytes := []byte{
		37, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		byte(Version1),                                 //Version1
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, //UUID
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
	slice, e := NewFromReader(reader)
	if e != nil {
		b.Error(e)
		return
	}
	writer := bytes.NewBuffer([]byte{})
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, e = slice.SerialiseWriter(writer)
		b.StopTimer()
		if e != nil {
			b.Error(e)
		}
		writer.Reset()
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
		37, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		byte(Version1),                                 //Version1
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, //UUID
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
		37, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		byte(Version1),                                 //Version1
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, //UUID
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
		37, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		byte(Version1),                                 //Version1
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, //UUID
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
		37, 0, 0, 0, 0, 0, 0, 0, //len(transaction)
		byte(Version1),                                 //Version1
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, //UUID
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
