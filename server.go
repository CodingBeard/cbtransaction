package cbtransaction

import (
	"errors"
	"fmt"
	"github.com/codingbeard/cbtransaction/transaction"
	"github.com/codingbeard/cbtransaction/transaction/cbslice"
	"github.com/codingbeard/cbutil"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var (
	uploadDir = "cbtransaction_upload"
	clientDir = "cbtransaction_client_buckets"
)

type transactionInsertQueueItem struct {
	action transaction.ActionEnum
	data   interface{}
}

type Server struct {
	storageProvider           Storage
	logger                    Logger
	errorHandler              ErrorHandler
	defaultEncryptionProvider Encryption
	encryptionProviders       []Encryption
	defaultEncodingProvider   Encoding
	encodingProviders         []Encoding
	client                    *Client
	master                    *Master
	dataDir                   string
	globalLock                *sync.Mutex
	transactionQueueLock      *sync.Mutex
	transactionInsertQueue    []transactionInsertQueueItem
}

type ServerConfig struct {
	Logger               Logger
	ErrorHandler         ErrorHandler
	DefaultEncryptionKey [8]byte
	EncryptionProviders  []Encryption
	DefaultEncodingKey   [8]byte
	EncodingProviders    []Encoding
	StorageProvider      Storage
	Concat               bool
	DestructiveCompact   bool
	DataDir              string
}

func NewServer(config ServerConfig) (*Server, error) {
	var defaultEncryptionProvider Encryption
	for _, provider := range config.EncryptionProviders {
		key := provider.GetKey()
		if key == config.DefaultEncryptionKey {
			defaultEncryptionProvider = provider
		}
	}
	if defaultEncryptionProvider == nil {
		return nil, errors.New("could not find default encryption provider")
	}
	var defaultEncodingProvider Encoding
	for _, provider := range config.EncodingProviders {
		key := provider.GetKey()
		if key == config.DefaultEncodingKey {
			defaultEncodingProvider = provider
		}
	}
	if defaultEncodingProvider == nil {
		return nil, errors.New("could not find default encoding provider")
	}
	client, e := NewClient(ClientConfig{
		Logger:               config.Logger,
		ErrorHandler:         config.ErrorHandler,
		DefaultEncryptionKey: defaultEncryptionProvider.GetKey(),
		EncryptionProviders:  config.EncryptionProviders,
		DefaultEncodingKey:   defaultEncodingProvider.GetKey(),
		EncodingProviders:    config.EncodingProviders,
		DataDir:              filepath.Join(config.DataDir, clientDir),
	})
	if e != nil {
		return nil, e
	}
	return &Server{
		logger:                    config.Logger,
		errorHandler:              config.ErrorHandler,
		defaultEncryptionProvider: defaultEncryptionProvider,
		encryptionProviders:       config.EncryptionProviders,
		defaultEncodingProvider:   defaultEncodingProvider,
		encodingProviders:         config.EncodingProviders,
		storageProvider:           config.StorageProvider,
		dataDir:                   config.DataDir,
		client:                    client,
		globalLock:                &sync.Mutex{},
		transactionQueueLock:      &sync.Mutex{},
	}, nil
}

func (s *Server) Run() error {
	e := s.client.Start()
	if e != nil {
		return e
	}

	e = os.MkdirAll(filepath.Join(s.dataDir, uploadDir), os.ModePerm)
	if e != nil {
		return e
	}

	cbutil.RepeatingTask{
		Sleep:      time.Second,
		SleepFirst: true,
		Run:        s.insertTransactionsQueue,
	}.Start()

	cbutil.RepeatingTask{
		Sleep:      time.Second,
		SleepFirst: true,
		Run:        s.negateExpiredTransactions,
	}.Start()

	cbutil.RepeatingTask{
		Sleep:      time.Minute,
		SleepFirst: true,
		Run:        s.uploadLatestBuckets,
	}.Start()

	select {}
}

func (s *Server) AddTransaction(action transaction.ActionEnum, data interface{}) {
	s.transactionQueueLock.Lock()
	defer s.transactionQueueLock.Unlock()
	s.transactionInsertQueue = append(
		s.transactionInsertQueue,
		transactionInsertQueueItem{action: action, data: data},
	)
}

func (s *Server) insertTransactionsQueue() {
	if len(s.transactionInsertQueue) == 0 {
		return
	}

	s.transactionQueueLock.Lock()

	insertItems := make([]transactionInsertQueueItem, len(s.transactionInsertQueue))
	copy(insertItems, s.transactionInsertQueue)
	s.transactionInsertQueue = s.transactionInsertQueue[:0]

	s.transactionQueueLock.Unlock()

	currentBucket := s.master.GetCurrentBucket()

	tempBucket, e := s.createTempBucketClone(currentBucket)
	if e != nil {
		s.transactionQueueLock.Lock()
		s.transactionInsertQueue = append(insertItems, s.transactionInsertQueue...)
		s.transactionQueueLock.Unlock()

		return
	}
	if tempBucket.GetFile() != nil {
		defer tempBucket.GetFile().Close()
	}
	tempBucket.Lock()

	for _, item := range insertItems {
		transaction := cbslice.NewVersion1()
		transactionId, e := uuid.NewUUID()
		if e != nil {
			s.errorHandler.Error(e)

			s.transactionQueueLock.Lock()
			s.transactionInsertQueue = append(insertItems, s.transactionInsertQueue...)
			s.transactionQueueLock.Unlock()
			_ = s.removeBucket(tempBucket)
			return
		}
		transaction.SetTransactionId(transactionId)
		transaction.SetActionEnum(item.action)
		transaction.SetEncryptionProviderKey(s.defaultEncryptionProvider.GetKey())
		transaction.SetEncodingProviderKey(s.defaultEncodingProvider.GetKey())
		encoded, e := s.defaultEncodingProvider.Encode(item.data)
		if e != nil {
			s.errorHandler.Error(e)

			s.transactionQueueLock.Lock()
			s.transactionInsertQueue = append(insertItems, s.transactionInsertQueue...)
			s.transactionQueueLock.Unlock()
			_ = s.removeBucket(tempBucket)
			return
		}
		transaction.SetData(s.defaultEncryptionProvider.Encrypt(encoded))
		_, e = transaction.SerialiseWriter(tempBucket.GetFile())
		if e != nil {
			s.errorHandler.Error(e)

			s.transactionQueueLock.Lock()
			s.transactionInsertQueue = append(insertItems, s.transactionInsertQueue...)
			s.transactionQueueLock.Unlock()
			_ = s.removeBucket(tempBucket)
			return
		}
	}

	tempBucket.Unlock()

	verifiedBucket, e := s.verifyBucket(tempBucket)
	if e != nil {
		s.transactionQueueLock.Lock()
		s.transactionInsertQueue = append(insertItems, s.transactionInsertQueue...)
		s.transactionQueueLock.Unlock()
		_ = s.removeBucket(tempBucket)
		return
	}

	if verifiedBucket.GetTransactionCount() != currentBucket.GetTransactionCount()+uint32(len(insertItems)) {
		s.errorHandler.Error(errors.New("transaction count of new bucket did not equal old count plus new items"))

		s.transactionQueueLock.Lock()
		s.transactionInsertQueue = append(insertItems, s.transactionInsertQueue...)
		s.transactionQueueLock.Unlock()
		_ = s.removeBucket(tempBucket)
		return
	}

	s.globalLock.Lock()
	defer s.globalLock.Unlock()

	currentBucket, e = s.replaceBucket(currentBucket, verifiedBucket)
	if e != nil {
		s.transactionQueueLock.Lock()
		s.transactionInsertQueue = append(insertItems, s.transactionInsertQueue...)
		s.transactionQueueLock.Unlock()
		_ = s.removeBucket(tempBucket)
		return
	}

	if verifiedBucket.GetFile() != nil {
		_ = verifiedBucket.GetFile().Close()
	}

	s.master.currentBucket = currentBucket
}

func (s *Server) createTempBucketClone(bucket *Bucket) (*Bucket, error) {
	tempFileName := bucket.GetFileName() + ".temp." + strconv.FormatInt(time.Now().UnixNano(), 10)
	tempFile, e := os.Create(filepath.Join(
		s.dataDir,
		tempFileName,
	))
	if e != nil {
		s.errorHandler.Error(e)
		return nil, e
	}

	bucket.Lock()
	defer bucket.Unlock()
	_, e = io.Copy(tempFile, bucket.GetFile())
	if e != nil {
		s.errorHandler.Error(e)
		return nil, e
	}

	tempBucket, e := NewBucketFromFile(tempFile)
	if e != nil {
		s.errorHandler.Error(e)
		return nil, e
	}

	tempBucket.SetVersion(bucket.GetVersion() + 1)
	tempBucket.SetFileName(tempFileName)
	tempBucket.SetModTime(time.Now().Unix())
	tempBucket.SetTransactionCount(bucket.GetTransactionCount())

	return tempBucket, nil
}

func (s *Server) removeBucket(bucket *Bucket) error {
	bucket.Lock()
	defer bucket.Unlock()
	return os.Remove(filepath.Join(s.dataDir, bucket.GetFileName()))
}

func (s *Server) verifyBucket(bucket *Bucket) (*Bucket, error) {
	bucket.Lock()
	defer bucket.Unlock()

	_, e := bucket.GetFile().Seek(0, io.SeekStart)
	if e != nil {
		s.errorHandler.Error(e)
		return nil, e
	}

	transactionCount := uint32(0)

	for true {
		oldPosition, posE := bucket.GetFile().Seek(0, io.SeekCurrent)
		if posE != nil {
			s.errorHandler.Error(posE)
			return nil, posE
		}

		_, e := cbslice.NewFromReader(bucket.GetFile())

		newPosition, posE := bucket.GetFile().Seek(0, io.SeekCurrent)
		if posE != nil {
			s.errorHandler.Error(posE)
			return nil, posE
		}

		if e != nil {
			if oldPosition != newPosition {
				err := fmt.Errorf("invalid transaction bucket: %s. %s", bucket.GetFileName(), e.Error())
				s.errorHandler.Error(err)
				return nil, err
			}
			break
		}

		transactionCount++
	}

	bucket.SetTransactionCount(transactionCount)

	return bucket, nil
}

func (s *Server) replaceBucket(bucket *Bucket, newBucket *Bucket) (*Bucket, error) {
	return nil, nil
}

func (s *Server) negateExpiredTransactions() {

}

func (s *Server) uploadLatestBuckets() {
	s.copyBucketsToUploadDir()

	s.concatUploadBuckets()

	s.generateUploadMaster()

	s.compressUploadBuckets()

	s.verifyUploadBuckets()

	s.uploadBuckets()
}

func (s *Server) copyBucketsToUploadDir() {
	s.globalLock.Lock()
	defer s.globalLock.Unlock()
}

func (s *Server) concatUploadBuckets() {

}

func (s *Server) generateUploadMaster() {

}

func (s *Server) compressUploadBuckets() {

}

func (s *Server) verifyUploadBuckets() {

}

func (s *Server) uploadBuckets() {

}
