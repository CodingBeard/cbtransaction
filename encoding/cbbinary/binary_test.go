package cbbinary

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name string
		args args
		want *Encoding
	}{
		{
			"little",
			args{config: Config{Endian: binary.LittleEndian}},
			&Encoding{endian: binary.LittleEndian},
		},
		{
			"big",
			args{config: Config{Endian: binary.BigEndian}},
			&Encoding{endian: binary.BigEndian},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncoding_Decode(t *testing.T) {
	type args struct {
		encoded []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
		out     int64
	}{
		{
			"int64",
			args{encoded: []byte{210, 4, 0, 0, 0, 0, 0, 0}},
			1234,
			false,
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(Config{Endian: binary.LittleEndian})
			err := b.Decode(tt.args.encoded, &tt.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.out, tt.want) {
				t.Errorf("Decode() gotN = %v, want %v", tt.out, tt.want)
			}
		})
	}
}

func TestEncoding_DecodeReader(t *testing.T) {
	type args struct {
		reader io.Reader
		out    int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"int64",
			args{
				reader: bytes.NewReader([]byte{210, 4, 0, 0, 0, 0, 0, 0}),
				out:    0,
			},
			1234,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(Config{Endian: binary.LittleEndian})
			err := b.DecodeReader(tt.args.reader, &tt.args.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.out != tt.want {
				t.Errorf("DecodeReader() out = %v, want %v", tt.args.out, tt.want)
			}
		})
	}
}

func TestEncoding_Encode(t *testing.T) {
	type args struct {
		data int64
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"int64",
			args{data: 1234},
			[]byte{210, 4, 0, 0, 0, 0, 0, 0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(Config{Endian: binary.LittleEndian})
			out, err := b.Encode(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(out, tt.want) {
				t.Errorf("Encode() gotN = %v, want %v", out, tt.want)
			}
		})
	}
}

func TestEncoding_EncodeWriter(t *testing.T) {
	type args struct {
		data int64
	}
	tests := []struct {
		name            string
		args            args
		wantWriterBytes []byte
		wantErr         bool
	}{
		{
			"int64",
			args{data: 1234},
			[]byte{210, 4, 0, 0, 0, 0, 0, 0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(Config{Endian: binary.LittleEndian})
			writer := &bytes.Buffer{}
			err := b.EncodeWriter(tt.args.data, writer)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.Bytes(); !reflect.DeepEqual(gotWriter, tt.wantWriterBytes) {
				t.Errorf("EncodeWriter() gotWriter = %v, want %v", gotWriter, tt.wantWriterBytes)
			}
		})
	}
}

func TestEncoding_GetKey(t *testing.T) {
	tests := []struct {
		name string
		want [8]byte
	}{
		{
			"little",
			[8]byte{'b', 'i', 'n', 'a', 'r', 'y', 'l', 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(Config{Endian: binary.LittleEndian})
			if got := b.GetKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
	tests2 := []struct {
		name string
		want [8]byte
	}{
		{
			"big",
			[8]byte{'b', 'i', 'n', 'a', 'r', 'y', 'b', 0},
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			b := New(Config{Endian: binary.BigEndian})
			if got := b.GetKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(Config{Endian:binary.LittleEndian})
	}
}

func BenchmarkEncoding_Encode1KB(b *testing.B) {
	b.StopTimer()
	encoder := New(Config{Endian:binary.LittleEndian})
	data := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		data[i] = 'a'
	}
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, e := encoder.Encode(data)
		b.StopTimer()
		if e != nil {
			b.Error(e)
			return
		}
	}
}

func BenchmarkEncoding_EncodeWriter1KB(b *testing.B) {
	b.StopTimer()
	encoder := New(Config{Endian:binary.LittleEndian})
	data := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		data[i] = 'a'
	}
	writer := bytes.NewBuffer(data)
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		e := encoder.EncodeWriter(data, writer)
		b.StopTimer()
		if e != nil {
			b.Error(e)
			return
		}
		writer.Reset()
	}
}

func BenchmarkEncoding_Decode1KB(b *testing.B) {
	b.StopTimer()
	encoder := New(Config{Endian:binary.LittleEndian})
	encoded := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		encoded[i] = 'a'
	}
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		var out [1000]byte
		e := encoder.Decode(encoded, &out)
		b.StopTimer()
		if e != nil {
			b.Error(e)
			return
		}
	}
}

func BenchmarkEncoding_DecodeReader1KB(b *testing.B) {
	b.StopTimer()
	encoder := New(Config{Endian:binary.LittleEndian})
	encoded := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		encoded[i] = 'a'
	}
	reader := bytes.NewReader(encoded)
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		var out [1000]byte
		e := encoder.DecodeReader(reader, &out)
		b.StopTimer()
		if e != nil {
			b.Error(e)
			return
		}
		_, _ = reader.Seek(0, 0)
	}
}