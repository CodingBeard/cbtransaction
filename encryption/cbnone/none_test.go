package cbnone

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestEncryption_Decrypt(t *testing.T) {
	type args struct {
		encrypted []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"string",
			args{encrypted: []byte("some string")},
			[]byte("some string"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Encryption{}
			if got := n.Decrypt(tt.args.encrypted); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncryption_DecryptReader(t *testing.T) {
	type args struct {
		reader io.Reader
		out    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"string",
			args{reader: bytes.NewReader([]byte("some string")), out: make([]byte, 11)},
			[]byte("some string"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Encryption{}
			if err := n.DecryptReader(tt.args.reader, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("DecryptReader() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.out, tt.want) {
				t.Errorf("Decrypt() = %v, want %v", tt.args.out, tt.want)
			}
		})
	}
}

func TestEncryption_Encrypt(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"string",
			args{data: []byte("some string")},
			[]byte("some string"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Encryption{}
			if got := n.Encrypt(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncryption_EncryptWriter(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name       string
		args       args
		wantWriter []byte
		wantErr    bool
	}{
		{
			"string",
			args{data: []byte("some string")},
			[]byte("some string"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Encryption{}
			writer := &bytes.Buffer{}
			err := n.EncryptWriter(tt.args.data, writer)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.Bytes(); !reflect.DeepEqual(gotWriter, tt.wantWriter) {
				t.Errorf("EncryptWriter() gotWriter = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestEncryption_GetKey(t *testing.T) {
	tests := []struct {
		name string
		want [8]byte
	}{
		{
			"default",
			[8]byte{'n', 'o', 'n', 'e', 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Encryption{}
			if got := n.GetKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
