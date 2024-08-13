package server

import (
	"fmt"
	"sync"
)

var ErrOffsetNotFound = fmt.Errorf("record at given offset not found")

type CommitLog struct {
	mu      sync.Mutex
	records []Record
}

func NewCommitLog() *CommitLog {
	return &CommitLog{}
}

type Record struct {
	Offset uint64 // This identifies the index at which the record will be stored in the log. We choose an unsigned integer here because we should never try to store a record at a negative number offset.
	Value  []byte
}

// Append adds a record onto the commit log.
func (cl *CommitLog) Append(record Record) (uint64, error) {
	// A mutex is used to ensure that only one
	// caller can access the commit log at a given time.
	cl.mu.Lock()
	defer cl.mu.Unlock()

	record.Offset = uint64(len(cl.records))
	cl.records = append(cl.records, record)

	return record.Offset, nil
}

// Read retrieves a record from our commit log. We specify the location in the
// commit log of our desired record using the offset that the record was originally written to.
func (cl *CommitLog) Read(offset uint64) (Record, error) {
	// A mutex is used to ensure that only one
	// caller can access the commit log at a given time.
	cl.mu.Lock()
	defer cl.mu.Unlock()

	if offset >= uint64(len(cl.records)) {
		return Record{}, ErrOffsetNotFound
	}

	return cl.records[offset], nil
}
