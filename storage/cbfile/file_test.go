package cbfile

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Storage
		wantErr bool
	}{
		{
			"emptyPath",
			args{config: Config{BasePath: ""}},
			nil,
			true,
		},
		{
			"invalidPath",
			args{config: Config{BasePath: "/a/b/c/"}},
			nil,
			true,
		},
		{
			"validPath",
			args{config: Config{BasePath: os.TempDir()}},
			&Storage{basePath: strings.TrimRight(os.TempDir(), "/")},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Concat(t *testing.T) {
	type fields struct {
		basePath string
	}
	type args struct {
		destination string
		filenames   []string
	}
	tempDir := os.TempDir()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func()
		test    func()
		cleanup func()
	}{
		{
			name:   "destinationError",
			fields: fields{basePath: tempDir},
			args: args{
				destination: ".",
			},
			wantErr: true,
			setup:   nil,
			cleanup: nil,
		},
		{
			name:   "emptySourceError",
			fields: fields{basePath: tempDir},
			args: args{
				destination: "destination.concat.emptySourceError.txt",
				filenames:   []string{},
			},
			wantErr: true,
			setup:   nil,
			cleanup: nil,
		},
		{
			name:   "invalidSourceError",
			fields: fields{basePath: tempDir},
			args: args{
				destination: "destination.concat.invalidSourceError.txt",
				filenames:   []string{"invalid.concat.invalidSourceError.txt"},
			},
			wantErr: true,
			setup:   nil,
			cleanup: nil,
		},
		{
			name:   "2files",
			fields: fields{basePath: tempDir},
			args: args{
				destination: "destination.concat.2files.txt",
				filenames:   []string{"source1.concat.2files.txt", "source2.concat.2files.txt"},
			},
			wantErr: false,
			setup: func() {
				e := ioutil.WriteFile(filepath.Join(tempDir, "source1.concat.2files.txt"), []byte("source1"), os.ModePerm)
				if e != nil {
					t.Error(e)
				}
				e = ioutil.WriteFile(filepath.Join(tempDir, "source2.concat.2files.txt"), []byte("source2"), os.ModePerm)
				if e != nil {
					t.Error(e)
				}
			},
			test: func() {
				content, e := ioutil.ReadFile(filepath.Join(tempDir, "destination.concat.2files.txt"))
				if e != nil {
					t.Error(e)
					return
				}
				expected := []byte("source1source2")
				if string(content) != string(expected) {
					t.Errorf("Concat() destination content = %s, want = %s", string(content), string(expected))
				}
			},
			cleanup: func() {
				e := os.Remove(filepath.Join(tempDir, "source1.concat.2files.txt"))
				if e != nil {
					t.Error(e)
				}
				e = os.Remove(filepath.Join(tempDir, "source2.concat.2files.txt"))
				if e != nil {
					t.Error(e)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				basePath: tt.fields.basePath,
			}
			if tt.setup != nil {
				tt.setup()
			}
			err := s.Concat(tt.args.destination, tt.args.filenames...)
			if (err != nil) != tt.wantErr {
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

func TestStorage_Delete(t *testing.T) {
	type fields struct {
		basePath string
	}
	type args struct {
		filename string
	}
	tempDir := os.TempDir()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func()
		test    func()
		cleanup func()
	}{
		{
			name:    "invalidPath",
			fields:  fields{basePath: tempDir},
			args:    args{filename: "."},
			wantErr: true,
			setup:   nil,
			test:    nil,
			cleanup: nil,
		},
		{
			name:    "validFile",
			fields:  fields{basePath: tempDir},
			args:    args{filename: "file.delete.validFile.txt"},
			wantErr: false,
			setup: func() {
				e := ioutil.WriteFile(filepath.Join(tempDir, "file.delete.validFile.txt"), []byte("file"), os.ModePerm)
				if e != nil {
					t.Error(e)
				}
			},
			test: func() {
				_, e := os.Stat(filepath.Join(tempDir, "file.delete.validFile.txt"))
				if e == nil {
					t.Errorf("Delete() error = %v, want not nil", e)
				}
			},
			cleanup: func() {
				_ = os.Remove(filepath.Join(tempDir, "file.delete.validFile.txt"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				basePath: tt.fields.basePath,
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

func TestStorage_Download(t *testing.T) {
	type fields struct {
		basePath string
	}
	type args struct {
		filename string
	}
	tempDir := os.TempDir()
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWriter string
		wantErr    bool
		setup      func()
		test       func()
		cleanup    func()
	}{
		{
			name:       "invalidFile",
			fields:     fields{basePath: tempDir},
			args:       args{filename: "."},
			wantWriter: "",
			wantErr:    true,
			setup:      nil,
			test:       nil,
			cleanup:    nil,
		},
		{
			name:       "fileDownload",
			fields:     fields{basePath: tempDir},
			args:       args{filename: "file.download.fileDownload.txt"},
			wantWriter: "file-contents",
			wantErr:    false,
			setup: func() {
				e := ioutil.WriteFile(filepath.Join(tempDir, "file.download.fileDownload.txt"), []byte("file-contents"), os.ModePerm)
				if e != nil {
					t.Error(e)
				}
			},
			cleanup: func() {
				_ = os.Remove(filepath.Join(tempDir, "file.download.fileDownload.txt"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			s := &Storage{
				basePath: tt.fields.basePath,
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

func TestStorage_Upload(t *testing.T) {
	type fields struct {
		basePath string
	}
	type args struct {
		filename string
		reader   io.Reader
	}
	tempDir := os.TempDir()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func()
		test    func()
		cleanup func()
	}{
		{
			name:    "invalidFile",
			fields:  fields{basePath: tempDir},
			args:    args{filename: "."},
			wantErr: true,
			setup:   nil,
			test:    nil,
			cleanup: nil,
		},
		{
			name:    "fileUpload",
			fields:  fields{basePath: tempDir},
			args:    args{filename: "file.download.fileUpload.txt", reader: bytes.NewBuffer([]byte("file-upload-contents"))},
			wantErr: false,
			test: func() {
				contents, e := ioutil.ReadFile(filepath.Join(tempDir, "file.download.fileUpload.txt"))
				if e != nil {
					t.Error(e)
					return
				}

				want := "file-upload-contents"
				if string(contents) != want {
					t.Errorf("Upload() uploaded file contents = %s, want = %s", string(contents), want)
				}
			},
			cleanup: func() {
				_ = os.Remove(filepath.Join(tempDir, "file.download.fileUpload.txt"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			s := &Storage{
				basePath: tt.fields.basePath,
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

func TestStorage_checkFilename(t *testing.T) {
	type fields struct {
		basePath string
	}
	type args struct {
		filename string
	}
	tempDir := os.TempDir()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func()
		test    func()
		cleanup func()
	}{
		{
			name:    "dot",
			fields:  fields{basePath: tempDir},
			args:    args{filename: "."},
			wantErr: true,
		},
		{
			name:    "dotdot",
			fields:  fields{basePath: tempDir},
			args:    args{filename: ".."},
			wantErr: true,
		},
		{
			name:    "outsideBasePath",
			fields:  fields{basePath: tempDir},
			args:    args{filename: "../file.txt"},
			wantErr: true,
		},
		{
			name:    "dir",
			fields:  fields{basePath: tempDir},
			args:    args{filename: "folder"},
			wantErr: true,
			setup: func() {
				e := os.Mkdir(filepath.Join(tempDir, "folder"), os.ModePerm)
				if e != nil {
					t.Error(e)
				}
			},
			cleanup: func() {
				e := os.Remove(filepath.Join(tempDir, "folder"))
				if e != nil {
					t.Error(e)
				}
			},
		},
		{
			name:    "validRelative",
			fields:  fields{basePath: tempDir},
			args:    args{filename: "folder/../file.txt"},
			wantErr: false,
			setup: func() {
				e := os.Mkdir(filepath.Join(tempDir, "folder"), os.ModePerm)
				if e != nil {
					t.Error(e)
				}
				e = ioutil.WriteFile(filepath.Join(tempDir, "file.txt"), []byte("file-contents"), os.ModePerm)
				if e != nil {
					t.Error(e)
				}
			},
			cleanup: func() {
				e := os.RemoveAll(filepath.Join(tempDir, "folder"))
				if e != nil {
					t.Error(e)
				}
			},
		},
		{
			name:    "validAbs",
			fields:  fields{basePath: tempDir},
			args:    args{filename: "folder/file.txt"},
			wantErr: false,
			setup: func() {
				e := os.Mkdir(filepath.Join(tempDir, "folder"), os.ModePerm)
				if e != nil {
					t.Error(e)
				}
				e = ioutil.WriteFile(filepath.Join(tempDir, "file.txt"), []byte("file-contents"), os.ModePerm)
				if e != nil {
					t.Error(e)
				}
			},
			cleanup: func() {
				e := os.RemoveAll(filepath.Join(tempDir, "folder"))
				if e != nil {
					t.Error(e)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			s := &Storage{
				basePath: tt.fields.basePath,
			}
			if err := s.checkFilename(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("checkFilename() error = %v, wantErr %v", err, tt.wantErr)
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
