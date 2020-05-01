package cbtransaction

import (
	"github.com/codingbeard/cbtransaction/encoding/cbmsgpack"
	"github.com/codingbeard/cbtransaction/encryption/cbnone"
	"github.com/codingbeard/cbtransaction/storage/cbfile"
	"github.com/codingbeard/cbtransaction/transaction"
	"os"
	"reflect"
	"testing"
)

var getDefaultTestServer = func() (*Server, error) {
	dir := os.TempDir()
	storage, e := cbfile.New(cbfile.Config{
		BasePath: dir,
	})
	if e != nil {
		return nil, e
	}
	return NewServer(ServerConfig{
		Logger:               defaultLogger{},
		ErrorHandler:         DefaultErrorHandler{},
		DefaultEncryptionKey: cbnone.Key,
		EncryptionProviders:  []Encryption{&cbnone.Encryption{}},
		DefaultEncodingKey:   cbmsgpack.Key,
		EncodingProviders:    []Encoding{cbmsgpack.New()},
		StorageProvider:      storage,
		Concat:               false,
		DestructiveCompact:   false,
		DataDir:              dir,
	})
}

func TestServer_AddTransaction(t *testing.T) {
	type args struct {
		action transaction.ActionEnum
		data   interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "add nil",
			args: args{
				action: transaction.ActionAdd,
				data:   nil,
			},
		},
		{
			name: "remove int",
			args: args{
				action: transaction.ActionRemove,
				data:   2147,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}

			s.AddTransaction(tt.args.action, tt.args.data)

			got := len(s.transactionInsertQueue)
			want := 1
			if got != want {
				t.Errorf("len(s.transactionInsertQueue) is incorrect, got %d, want %d", got, want)
				return
			}

			gotAction := s.transactionInsertQueue[0].action
			if !reflect.DeepEqual(gotAction, tt.args.action) {
				t.Errorf("s.transactionInsertQueue[0].action is incorrect, got %v, want %v", gotAction, tt.args.action)
			}

			gotData := s.transactionInsertQueue[0].data
			if !reflect.DeepEqual(gotData, tt.args.data) {
				t.Errorf("s.transactionInsertQueue[0].data is incorrect, got %v, want %v", gotData, tt.args.data)
			}
		})
	}
}

func TestServer_Run(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			if err := s.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.cleanup()
		})
	}
}

func TestServer_compressUploadBuckets(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			tt.cleanup()
		})
	}
}

func TestServer_concatUploadBuckets(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			tt.cleanup()
		})
	}
}

func TestServer_copyBucketsToUploadDir(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			tt.cleanup()
		})
	}
}

func TestServer_createTempBucketClone(t *testing.T) {
	type args struct {
		bucket *Bucket
	}
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
		args    args
		want    *Bucket
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			got, err := s.createTempBucketClone(tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("createTempBucketClone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createTempBucketClone() got = %v, want %v", got, tt.want)
			}
			tt.cleanup()
		})
	}
}

func TestServer_generateUploadMaster(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			tt.cleanup()
		})
	}
}

func TestServer_insertTransactionsQueue(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			tt.cleanup()
		})
	}
}

func TestServer_negateExpiredTransactions(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			tt.cleanup()
		})
	}
}

func TestServer_removeBucket(t *testing.T) {
	type args struct {
		bucket *Bucket
	}
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			if err := s.removeBucket(tt.args.bucket); (err != nil) != tt.wantErr {
				t.Errorf("removeBucket() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.cleanup()
		})
	}
}

func TestServer_replaceBucket(t *testing.T) {
	type args struct {
		bucket    *Bucket
		newBucket *Bucket
	}
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
		args    args
		want    *Bucket
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			got, err := s.replaceBucket(tt.args.bucket, tt.args.newBucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("replaceBucket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("replaceBucket() got = %v, want %v", got, tt.want)
			}
			tt.cleanup()
		})
	}
}

func TestServer_uploadBuckets(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			tt.cleanup()
		})
	}
}

func TestServer_uploadLatestBuckets(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			tt.cleanup()
		})
	}
}

func TestServer_verifyBucket(t *testing.T) {
	type args struct {
		bucket *Bucket
	}
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
		args    args
		want    *Bucket
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			got, err := s.verifyBucket(tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyBucket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("verifyBucket() got = %v, want %v", got, tt.want)
			}
			tt.cleanup()
		})
	}
}

func TestServer_verifyUploadBuckets(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(s *Server) *Server
		cleanup func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := getDefaultTestServer()
			if e != nil {
				t.Error(e.Error())
				return
			}
			s = tt.setup(s)
			tt.cleanup()
		})
	}
}
