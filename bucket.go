package cbtransaction

import (
	"io"
	"sync"
)

type Bucket struct {
	FileName        string
	Hash            string
	CompressedHash  string
	CompressionAlgo string
	ModTime         int64
	Version         uint32

	// used internally / not persisted
	lock *sync.RWMutex
	file io.ReadWriteSeeker
}

func NewBucketFromFile(file io.ReadWriteSeeker) (*Bucket, error) {
	return &Bucket{
		lock: &sync.RWMutex{},
		file: file,
	}, nil
}

func (b *Bucket) GetFile() io.ReadWriteSeeker {
	return b.file
}

func (b *Bucket) Lock() {
	b.lock.Lock()
}

func (b *Bucket) Unlock() {
	b.lock.Unlock()
}

func (b *Bucket) GetFileName() string {
	return b.FileName
}

func (b *Bucket) SetFileName(fileName string) {
	b.FileName = fileName
}

func (b *Bucket) GetHash() string {
	return b.Hash
}

func (b *Bucket) SetHash(hash string) {
	b.Hash = hash
}

func (b *Bucket) GetCompressedHash() string {
	return b.CompressedHash
}

func (b *Bucket) SetCompressedHash(compressedHash string) {
	b.CompressedHash = compressedHash
}

func (b *Bucket) GetCompressionAlgo() string {
	return b.CompressionAlgo
}

func (b *Bucket) SetCompressionAlgo(compressionAlgo string) {
	b.CompressionAlgo = compressionAlgo
}

func (b *Bucket) GetModTime() int64 {
	return b.ModTime
}

func (b *Bucket) SetModTime(modTime int64) {
	b.ModTime = modTime
}

func (b *Bucket) GetVersion() uint32 {
	return b.Version
}

func (b *Bucket) SetVersion(version uint32) {
	b.Version = version
}
