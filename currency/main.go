package main

import (
	"log"
	"net"
	"os"

	"github.com/danielhoward314/microservices-test/currency/data"
	protos "github.com/danielhoward314/microservices-test/currency/protos"
	"github.com/danielhoward314/microservices-test/currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// instantiate dependencies of server struct
	l := log.New(os.Stdout, "currency", log.LstdFlags)
	er, err := data.NewRates()
	if err != nil {
		l.Print("unable to get exchange rates")
		os.Exit(1)
	}
	// instantiate a struct representing the generated server interface
	cs := server.NewCurrency(l, er)
	// instantiate a grpc server
	gs := grpc.NewServer()
	// bind the server interface to the grpc server (kind of like registering a REST handler)
	protos.RegisterCurrencyServer(gs, cs)
	// to enable grpcurl to ask grpc server to list available services and methods (not for PRD)
	reflection.Register(gs)
	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		l.Print("unable to listen on tcp port")
		os.Exit(1)
	}
	gs.Serve(lis)
}
