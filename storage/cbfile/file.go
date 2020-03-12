package cbfile

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Storage struct {
	basePath string
}

type Config struct {
	BasePath string
}

func New(config Config) (*Storage, error) {
	basePath := strings.TrimRight(config.BasePath, "/")
	stat, e := os.Stat(basePath)
	if e != nil {
		return nil, e
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("basePath is not a valid directory: %s", basePath)
	}
	return &Storage{
		basePath: basePath,
	}, nil
}

func (s *Storage) checkFilename(filename string) error {
	if filename == "" || filename == "." || filename == ".." {
		return fmt.Errorf("invalid filename: %s", filename)
	}

	abs, e := filepath.Abs(path.Join(s.basePath, filename))
	if e != nil {
		return e
	}

	if !strings.HasPrefix(abs, s.basePath) {
		return fmt.Errorf("absolute file path (%s) is outside of base path (%s)", abs, s.basePath)
	}

	stat, e := os.Stat(abs)
	if e == nil {
		if stat.IsDir() {
			return fmt.Errorf("file is a directory: %s", abs)
		}
	}

	return nil
}

func (s *Storage) Upload(filename string, reader io.Reader) error {
	if e := s.checkFilename(filename); e != nil {
		return e
	}
	writer, e := os.Create(path.Join(s.basePath, filename))
	if e != nil {
		return e
	}
	if writer != nil {
		defer writer.Close()
	}
	_, e = io.Copy(writer, reader)
	return e
}

func (s *Storage) Download(filename string, writer io.Writer) error {
	if e := s.checkFilename(filename); e != nil {
		return e
	}
	reader, e := os.Open(path.Join(s.basePath, filename))
	if e != nil {
		return e
	}
	if reader != nil {
		defer reader.Close()
	}
	_, e = io.Copy(writer, reader)
	return e
}

func (s *Storage) Concat(destination string, filenames ...string) error {
	if e := s.checkFilename(destination); e != nil {
		return e
	}
	if len(filenames) == 0 {
		return errors.New("no filenames given to concat")
	}
	for _, filename := range filenames {
		if e := s.checkFilename(filename); e != nil {
			return e
		}
		filePath := path.Join(s.basePath, filename)
		stat, e := os.Stat(filePath)
		if e != nil {
			return e
		}
		if stat.IsDir() {
			return fmt.Errorf("cannot concat file into destination, %s is a directory", filePath)
		}
	}
	writer, e := os.Create(path.Join(s.basePath, destination))
	if e != nil {
		return e
	}
	if writer != nil {
		defer writer.Close()
	}
	for _, filename := range filenames {
		reader, e := os.Open(path.Join(s.basePath, filename))
		if e != nil {
			return e
		}
		_, e = io.Copy(writer, reader)
		_ = reader.Close()
		if e != nil {
			return e
		}
	}

	return nil
}

func (s *Storage) Delete(filename string) error {
	return os.Remove(path.Join(s.basePath, filename))
}
