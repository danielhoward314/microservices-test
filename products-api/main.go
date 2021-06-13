package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	protos "github.com/danielhoward314/microservices-test/currency/protos"
	"github.com/danielhoward314/microservices-test/products-api/handlers"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// instantiate mux router
	r := mux.NewRouter()
	// instantiate dependencies to be injected into handler structs
	l := log.New(os.Stdout, "products-api", log.LstdFlags)
	// instantiate grpc client
	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	cc := protos.NewCurrencyClient(conn)
	// instantiate handler
	ph := handlers.NewProducts(l, cc)
	// create a subrouter for all routes assoc'd w/ products handler
	psr := r.PathPrefix("/products").Subrouter()
	// create subrouter routes that bind paths to methods of products handler
	psr.HandleFunc("", ph.AddProduct).Methods(http.MethodPost)
	psr.HandleFunc("", ph.GetProducts).Methods(http.MethodGet)
	psr.HandleFunc("/{id:[0-9]+}", ph.GetProduct).Methods(http.MethodGet)
	/*
		gorilla/mux path params syntax with regex, becomes accessible in handler through:
		```
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		```
	*/
	psr.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct).Methods(http.MethodPut)
	// CORS handler (assuming a React app served on port 3000 that makes XHR requests to this API server on 8080)
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))
	// create a server with configs (note the CORS handler wrapping the mux router)
	s := &http.Server{
		Addr:         ":8080",
		Handler:      ch(r),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	// wrap listenAndServe call in a goroutine, since it blocks
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server %s\n", err)
			os.Exit(1)
		}
	}()
	// create a buffered channel for OS signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	/*
		receiving on a channel blocks, which is what allows the
		goroutine above to run the server indefinitely
		until there is a signal from the OS to stop the process
	*/
	sig := <-c
	l.Printf("Got an os.Signal: %v", sig)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	s.Shutdown(ctx)
}
