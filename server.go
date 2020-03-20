package cbtransaction

import (
	"errors"
	"github.com/codingbeard/cbtransaction/transaction/cbslice"
	"github.com/codingbeard/cbutil"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	uploadDir = "cbtransaction_upload"
	clientDir = "cbtransaction_client_buckets"
)

type transactionInsertQueueItem struct {
	action ActionEnum
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
		Sleep:      time.Millisecond * 100,
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

func (s *Server) AddTransaction(action ActionEnum, data interface{}) {
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

	s.globalLock.Lock()
	defer s.globalLock.Unlock()

	currentBucket := s.master.GetCurrentBucket()

	newBucket, e := s.duplicateBucket(currentBucket)
	if e != nil {
		s.errorHandler.Error(e)

		s.transactionQueueLock.Lock()
		s.transactionInsertQueue = append(insertItems, s.transactionInsertQueue...)
		s.transactionQueueLock.Unlock()

		return
	}

	for _, item := range insertItems {
		transaction := cbslice.NewVersion1()
		transactionId, e := uuid.NewUUID()
		if e != nil {
			s.errorHandler.Error(e)

			s.transactionQueueLock.Lock()
			s.transactionInsertQueue = append(insertItems, s.transactionInsertQueue...)
			s.transactionQueueLock.Unlock()
			// todo remove duplicated bucket
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
			// todo remove duplicated bucket
			return
		}
		transaction.SetData(s.defaultEncryptionProvider.Encrypt(encoded))
		_, e = transaction.SerialiseWriter(newBucket.GetFile())
		if e != nil {
			s.errorHandler.Error(e)

			s.transactionQueueLock.Lock()
			s.transactionInsertQueue = append(insertItems, s.transactionInsertQueue...)
			s.transactionQueueLock.Unlock()
			// todo remove duplicated bucket
			return
		}
	}

	e = s.replaceBucket(currentBucket, newBucket)
	if e != nil {
		s.errorHandler.Error(e)

		s.transactionQueueLock.Lock()
		s.transactionInsertQueue = append(insertItems, s.transactionInsertQueue...)
		s.transactionQueueLock.Unlock()
		// todo remove duplicated bucket
		return
	}
}

func (s *Server) duplicateBucket(bucket *Bucket) (*Bucket, error) {
	return nil, nil
}

func (s *Server) replaceBucket(bucket *Bucket, newBucket *Bucket) error {
	return nil
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
