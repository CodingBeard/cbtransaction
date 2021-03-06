package cbtransaction

import (
	"sync"
)

type Master struct {
	Buckets []Bucket

	// used internally / not persisted
	lock          *sync.RWMutex
	file          ReadWriteSeekCloser
	buckets       []*Bucket
	currentBucket *Bucket
}

func NewMasterFromFile(file ReadWriteSeekCloser) (*Master, error) {
	return &Master{
		lock: &sync.RWMutex{},
		file: file,
	}, nil
}

func (m *Master) GetFile() ReadWriteSeekCloser {
	return m.file
}

func (m *Master) Lock() {
	m.lock.Lock()
}

func (m *Master) Unlock() {
	m.lock.Unlock()
}

func (m *Master) GetBuckets() []*Bucket {
	return m.buckets
}

func (m *Master) GetCurrentBucket() *Bucket {
	return m.currentBucket
}

func (m *Master) SaveBucket(bucket *Bucket) {
	for key := range m.buckets {
		if m.buckets[key].GetFileName() == bucket.GetFileName() {
			*m.buckets[key] = *bucket
			return
		}
	}

	m.buckets = append(m.buckets, bucket)
}
