package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// struct type that implements the Handler interface
// sort of like a constructor defining dependencies
type Hello struct {
	l *log.Logger
}

// NewHello creates a new hello handler with the given logger
// example of DI in Go
func NewHello(l *log.Logger) *Hello {
	// syntax for instantiating a struct
	return &Hello{l}
}

// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle Hello requests")

	// read the body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.l.Println("Error reading body", err)

		http.Error(rw, "Unable to read request body", http.StatusBadRequest)
		return
	}

	// write the response
	fmt.Fprintf(rw, "Hello %s", b)
}
