package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielhoward314/microservices-test/product-images/files"
	"github.com/danielhoward314/microservices-test/product-images/handlers"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	basePath = "./imagestore"
)

func main() {
	r := mux.NewRouter()
	l := log.New(os.Stdout, "product-images", log.LstdFlags)
	store, err := files.NewLocal(1024*1000*5, basePath)
	if err != nil {
		l.Print("unable to create image store")
		os.Exit(1)
	}
	ih := handlers.NewImages(l, basePath, store)
	mw := handlers.GzipHandler{}
	ir := r.PathPrefix("/api/v1").Subrouter()
	ir.HandleFunc("/images/{filename:[a-zA-Z0-9\\-\\_]+\\.(?:gif|jpe?g|tiff?|png|webp|bmp)$}", ih.Upload).Methods(http.MethodPost)
	ir.HandleFunc("/images/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}/{filename:[a-zA-Z0-9\\-\\_]+\\.(?:gif|jpe?g|tiff?|png|webp|bmp)$}",
		ih.GetImage).Methods(http.MethodGet)
	ir.Use(mw.GzipMiddleware)
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))
	s := &http.Server{
		Addr:         ":9091",
		Handler:      ch(r),
		IdleTimeout:  time.Second * 120,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Print("could not start server")
			os.Exit(1)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c
	l.Printf("received interrupt or sigterm signal: %v", sig)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	s.Shutdown(ctx)
}
