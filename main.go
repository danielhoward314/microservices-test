package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielhoward314/microservices-test/handlers"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "test-api", log.LstdFlags)
	// instantiate the router
	r := mux.NewRouter()
	// instantiate the handlers
	productsHandler := handlers.NewProducts(l)
	// register the handler with a router
	pSub := r.PathPrefix("/products").Subrouter()
	pSub.HandleFunc("", productsHandler.GetProducts).Methods("GET")
	pSub.HandleFunc("", productsHandler.AddProduct).Methods("POST")
	pSub.HandleFunc("/{id:[0-9]+}", productsHandler.UpdateProduct).Methods("PUT")

	// rather than just `http.ListenAndServe(":8080", sm)`
	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// since ListenAndServe blocks, run inside go routine
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// whenever OS interrupt or kill command is received
	// the signal is sent on the channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	// blocks until os.Signal data is consumed from channel
	sig := <-c
	l.Println("Received interrupt or kill signal, graceful shutdown", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
