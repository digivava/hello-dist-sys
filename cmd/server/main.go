package main

import (
	"fmt"
	"log"

	"github.com/digivava/hello-dist-sys/internal/server"
)

func main() {
	port := "8080"
	log.Printf("Serving traffic on port %s...", port)
	srv := server.NewHTTPServer(fmt.Sprintf(":%s", port))
	log.Fatal(srv.ListenAndServe())
}
