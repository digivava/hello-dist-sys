package server

import (
	"fmt"
	"sync"
)

var ErrOffsetNotFound = fmt.Errorf("record at given offset not found")

// A WriteAheadLog is a data structure for storing a list of records,
// where new records can be appended to it, but old records cannot be modified later.
// This is good for use cases where order and accurate history is crucial, like a transaction ledger.
type WriteAheadLog struct {
	mu      sync.Mutex
	records []Record
}

func NewWriteAheadLog() *WriteAheadLog {
	return &WriteAheadLog{}
}

type Record struct {
	Offset uint64 `json:"offset"` // This identifies the index at which the record will be stored in the log. We choose an unsigned integer here because we should never try to store a record at a negative number offset.
	Value  []byte `json:"value"`
}

// Append adds a record onto the write-ahead log.
func (wal *WriteAheadLog) Append(record Record) (uint64, error) {
	// A mutex is used to ensure that only one
	// caller can access the write-ahead log at a given time.
	wal.mu.Lock()
	defer wal.mu.Unlock()

	record.Offset = uint64(len(wal.records))
	wal.records = append(wal.records, record)

	return record.Offset, nil
}

// Read retrieves a record from our write-ahead log. We specify the location in the
// log of our desired record using the offset that the record was originally written to.
func (wal *WriteAheadLog) Read(offset uint64) (Record, error) {
	// A mutex is used to ensure that only one
	// caller can access the write-ahead log at a given time.
	wal.mu.Lock()
	defer wal.mu.Unlock()

	if offset >= uint64(len(wal.records)) {
		return Record{}, ErrOffsetNotFound
	}

	return wal.records[offset], nil
}
