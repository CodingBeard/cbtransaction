package cbgooglebucket

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/option"
	"io"
)

type Storage struct {
	client *storage.Client
	bucket string
}

type Config struct {
	Bucket          string
	CredentialsJson []byte
}

func New(config Config) (*Storage, error) {
	if config.Bucket == "" {
		return nil, errors.New("invalid bucket name")
	}
	if len(config.CredentialsJson) == 0 {
		return nil, errors.New("empty credentials json")
	}
	client, e := storage.NewClient(context.Background(), option.WithCredentialsJSON(config.CredentialsJson))
	if e != nil {
		return nil, e
	}

	return &Storage{bucket: config.Bucket, client: client}, nil
}

func (s *Storage) checkFilename(filename string) error {
	if filename == "" || filename == "." || filename == ".." {
		return fmt.Errorf("invalid filename: %s", filename)
	}

	return nil
}

func (s *Storage) Upload(filename string, reader io.Reader) error {
	if e := s.checkFilename(filename); e != nil {
		return e
	}
	object := s.client.Bucket(s.bucket).Object(filename)

	ctx := context.Background()
	writer := object.NewWriter(ctx)

	_, e := io.Copy(writer, reader)
	if e != nil {
		return e
	}

	if e := writer.Close(); e != nil {
		return e
	}

	return nil
}

func (s *Storage) Download(filename string, writer io.Writer) error {
	if e := s.checkFilename(filename); e != nil {
		return e
	}

	object := s.client.Bucket(s.bucket).Object(filename)

	ctx := context.Background()
	reader, e := object.NewReader(ctx)
	if e != nil {
		return e
	}

	_, e = io.Copy(writer, reader)
	if e != nil {
		return e
	}

	if e := reader.Close(); e != nil {
		return e
	}

	return nil
}

func (s *Storage) Concat(destination string, filenames ...string) error {
	if e := s.checkFilename(destination); e != nil {
		return e
	}

	dst := s.client.Bucket(s.bucket).Object(destination)

	var objects []*storage.ObjectHandle
	for _, filename := range filenames {
		objects = append(objects, s.client.Bucket(s.bucket).Object(filename))
	}

	_, e := dst.ComposerFrom(objects...).Run(context.Background())

	return e
}

func (s *Storage) Delete(filename string) error {
	if e := s.checkFilename(filename); e != nil {
		return e
	}
	object := s.client.Bucket(s.bucket).Object(filename)
	return object.Delete(context.Background())
}
