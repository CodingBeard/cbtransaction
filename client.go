package cbtransaction

import (
	"errors"
)

type Client struct {
	logger                    Logger
	errorHandler              ErrorHandler
	encryptionProviders       []Encryption
	defaultEncryptionProvider Encryption
	encodingProviders         []Encoding
	defaultEncodingProvider   Encoding
	dataDir                   string
	master                    *Master
}

type ClientConfig struct {
	Logger               Logger
	ErrorHandler         ErrorHandler
	DefaultEncryptionKey [8]byte
	EncryptionProviders  []Encryption
	DefaultEncodingKey   [8]byte
	EncodingProviders    []Encoding
	DataDir              string
}

func NewClient(config ClientConfig) (*Client, error) {
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
	return &Client{
		logger:                    config.Logger,
		errorHandler:              config.ErrorHandler,
		defaultEncryptionProvider: defaultEncryptionProvider,
		encryptionProviders:       config.EncryptionProviders,
		defaultEncodingProvider:   defaultEncodingProvider,
		encodingProviders:         config.EncodingProviders,
		dataDir:                   config.DataDir,
	}, nil
}

func (s *Client) Start() error {
	return nil
}

func (s *Client) Verify() (bool, error) {
	return true, nil
}

func (s *Client) Download() error {
	return nil
}

func (s *Client) GetTransactions(currentVersion uint64, limit uint64) []Transaction {
	var transactions []Transaction

	return transactions
}

func (s *Client) Decrypt(transaction Transaction) (Transaction, error) {
	return transaction, nil
}

func (s *Client) Decode(transaction Transaction) (interface{}, error) {
	return transaction, nil
}

func (s *Client) GetMaster() *Master {
	return s.master
}
