# hello-dist-sys

This sample application is meant to be used as a resource for learning about distributed systems concepts in Go.

The code is based on the content in Travis Jeffery's book [Distributed Services with Go](https://pragprog.com/titles/tjgo/distributed-services-with-go/), with modifications and extra editorial comments to make things more comprehensible for developers new to distributed systems concepts. 

## Concepts

This application is an example of a service that manages a **write-ahead log** (WAL). A write-ahead log (also known as a "commit log" or "transaction log") is a data structure for storing a list of records, where new records can be appended to it, but old records cannot be modified later. This is important for use cases where order and accurate history is crucial, like transaction ledgers, storage engines, version control systems, or consensus algorithms.

For other concepts used in this repo, click the items in the list below to find a README with more information.

* [Protocol Buffers](./api/v1/README.md): We use this format for our internal service communications.

## Running the application

To start up the application's web server so you can start making HTTP requests against it, run `make run`.

Then you can create new records with your HTTP client of choice like this: 

`curl -XPOST http://localhost:8080 -d '{"record": {"value": "YmxhcmdoCg=="}}'`. 

This will return the offset for your newly created record in the write-ahead log. 

(We pass a base64-encoded value because we're working with the `[]byte` type in Go, and Go's `encoding/json` package encodes `[]byte` as a base64-encoded string. Use `echo "foo" | base64` to create a base64-encoded string.)

To fetch an existing record, use `curl -XGET http://localhost:8080/{offset}`, where offset is the unique location in the write-ahead log that was returned to you when you created a record. It should be a number, like 0 for the first record you create.

To compile new generated code if you change the protobuf definitions, you can run `make compile`.