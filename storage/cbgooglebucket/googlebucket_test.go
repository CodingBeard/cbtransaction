package cbgooglebucket

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"testing"
)

var bucket = flag.String("bucket", "", "bucket name")
var credentialsJsonPath = flag.String("credentialsJsonPath", "", "path to google json credentials for the bucket")

func getStorage() (*Storage, error) {
	if *bucket == "" || *credentialsJsonPath == "" {
		return nil, errors.New("please provide a bucket and the path to the json credentials: -bucket=mybucket -credentialsJsonPath=../../../creds.json")
	}

	contents, e := ioutil.ReadFile(*credentialsJsonPath)
	if e != nil {
		return nil, e
	}

	storage, e := New(Config{
		Bucket:          *bucket,
		CredentialsJson: contents,
	})
	if e != nil {
		return nil, e
	}

	return storage, nil
}

func TestStorage_Upload(t *testing.T) {
	type args struct {
		filename string
		reader   io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func()
		test    func()
		cleanup func()
	}{
		{
			name:    "invalidFilename",
			args:    args{filename: ".", reader: bytes.NewReader([]byte{})},
			wantErr: true,
		},
		{
			name:    "emptyReader",
			args:    args{filename: "cbtransaction-test/upload.emptyReader.txt", reader: bytes.NewReader([]byte{})},
			wantErr: false,
			cleanup: func() {
				s, e := getStorage()
				if e != nil {
					t.Error(e)
					return
				}
				_ = s.Delete("cbtransaction-test/upload.emptyReader.txt")
			},
		},
		{
			name:    "filledReader",
			args:    args{filename: "cbtransaction-test/upload.filledReader.txt", reader: bytes.NewReader([]byte("contents"))},
			wantErr: false,
			cleanup: func() {
				s, e := getStorage()
				if e != nil {
					t.Error(e)
					return
				}
				_ = s.Delete("cbtransaction-test/upload.filledReader.txt")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getStorage()
			if e != nil {
				t.Error(e)
				return
			}
			if tt.setup != nil {
				tt.setup()
			}
			if err := s.Upload(tt.args.filename, tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.test != nil {
				tt.test()
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestStorage_Download(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
		setup      func()
		test       func()
		cleanup    func()
	}{
		{
			name:    "invalidFilename",
			args:    args{filename: "."},
			wantErr: true,
		},
		{
			name:    "nonExistentFile",
			args:    args{filename: "asdfasdfilhbwrgilsdfviefv.txt"},
			wantErr: true,
		},
		{
			name:       "validFile",
			args:       args{filename: "cbtransaction-test/download.validFile.txt"},
			wantErr:    false,
			wantWriter: "content",
			setup: func() {
				s, e := getStorage()
				if e != nil {
					t.Error(e)
					return
				}
				e = s.Upload("cbtransaction-test/download.validFile.txt", bytes.NewReader([]byte("content")))
				if e != nil {
					t.Error(e)
				}
			},
			cleanup: func() {
				s, e := getStorage()
				if e != nil {
					t.Error(e)
					return
				}
				_ = s.Delete("cbtransaction-test/download.validFile.txt")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getStorage()
			if e != nil {
				t.Error(e)
				return
			}
			if tt.setup != nil {
				tt.setup()
			}
			writer := &bytes.Buffer{}
			err := s.Download(tt.args.filename, writer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Download() gotWriter = %v, want %v", gotWriter, tt.wantWriter)
			}
			if tt.test != nil {
				tt.test()
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func()
		test    func()
		cleanup func()
	}{
		{
			name:    "invalidFilename",
			args:    args{filename: "."},
			wantErr: true,
		},
		{
			name:    "validFile",
			args:    args{filename: "cbtransaction-test/delete.validFile.txt"},
			wantErr: false,
			setup: func() {
				s, e := getStorage()
				if e != nil {
					t.Error(e)
					return
				}
				e = s.Upload("cbtransaction-test/delete.validFile.txt", bytes.NewReader([]byte("content")))
				if e != nil {
					t.Error(e)
				}
			},
			test: func() {
				s, e := getStorage()
				if e != nil {
					t.Error(e)
					return
				}
				writer := bytes.NewBuffer([]byte{})
				e = s.Download("cbtransaction-test/delete.validFile.txt", writer)
				if e == nil {
					t.Error("Delete() error = nil, wantErr 404 file does not exist")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getStorage()
			if e != nil {
				t.Error(e)
				return
			}
			if tt.setup != nil {
				tt.setup()
			}
			if err := s.Delete(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.test != nil {
				tt.test()
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

func TestStorage_Concat(t *testing.T) {
	type args struct {
		destination string
		filenames   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func()
		test    func()
		cleanup func()
	}{
		{
			name:    "invalidFilename",
			args:    args{destination: "."},
			wantErr: true,
		},
		{
			name: "missingFile",
			args: args{
				destination: "cbtransaction-test/concat.missingFile.txt",
				filenames: []string{
					"cbtransaction-test/concat1.missingFile.txt",
					"cbtransaction-test/concat2.missingFile.txt",
				},
			},
			wantErr: true,
			setup: func() {
				s, e := getStorage()
				if e != nil {
					t.Error(e)
					return
				}
				e = s.Upload("cbtransaction-test/concat1.missingFile.txt", bytes.NewReader([]byte("content")))
				if e != nil {
					t.Error(e)
				}
			},
			cleanup: func() {
				s, e := getStorage()
				if e != nil {
					t.Error(e)
					return
				}
				_ = s.Delete("cbtransaction-test/concat.missingFile.txt")
				_ = s.Delete("cbtransaction-test/concat1.missingFile.txt")
			},
		},
		{
			name: "validFiles",
			args: args{
				destination: "cbtransaction-test/concat.validFiles.txt",
				filenames: []string{
					"cbtransaction-test/concat1.validFiles.txt",
					"cbtransaction-test/concat2.validFiles.txt",
				},
			},
			wantErr: false,
			setup: func() {
				s, e := getStorage()
				if e != nil {
					t.Error(e)
					return
				}
				e = s.Upload("cbtransaction-test/concat1.validFiles.txt", bytes.NewReader([]byte("content1")))
				if e != nil {
					t.Error(e)
				}
				e = s.Upload("cbtransaction-test/concat2.validFiles.txt", bytes.NewReader([]byte("content2")))
				if e != nil {
					t.Error(e)
				}
			},
			test: func() {
				s, e := getStorage()
				if e != nil {
					t.Error(e)
					return
				}
				writer := bytes.NewBuffer([]byte{})
				e = s.Download("cbtransaction-test/concat.validFiles.txt", writer)
				if e != nil {
					t.Error(e)
					return
				}
				expected := "content1content2"
				if writer.String() != expected {
					t.Errorf("Concat() content of destination file = %s, want = %s", writer.String(), expected)
				}
			},
			cleanup: func() {
				s, e := getStorage()
				if e != nil {
					t.Error(e)
					return
				}
				_ = s.Delete("cbtransaction-test/concat.validFiles.txt")
				_ = s.Delete("cbtransaction-test/concat1.validFiles.txt")
				_ = s.Delete("cbtransaction-test/concat2.validFiles.txt")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getStorage()
			if e != nil {
				t.Error(e)
				return
			}
			if tt.setup != nil {
				tt.setup()
			}
			if err := s.Concat(tt.args.destination, tt.args.filenames...); (err != nil) != tt.wantErr {
				t.Errorf("Concat() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.test != nil {
				tt.test()
			}
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}
