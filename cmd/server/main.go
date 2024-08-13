package main

import (
	"log"
	"os"

	"github.com/digivava/hello-dist-sys/internal/server"
)

func main() {
	cl := server.NewCommitLog()
	record := server.Record{Value: []byte("blah")}
	offset, err := cl.Append(record)
	if err != nil {
		log.Printf("unable to append record to log: %v", err)
		os.Exit(1)
	}
	log.Printf("added record with offset %d to log", offset)
	retrieved, err := cl.Read(offset)
	if err != nil {
		log.Printf("unable to read record from log: %v", err)
		os.Exit(1)
	}
	log.Printf("retrieved record with offset %d from log, content is \"%s\"", retrieved.Offset, retrieved.Value)
}
