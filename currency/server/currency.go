package server

import (
	"context"
	"log"

	protos "github.com/danielhoward314/microservices-test/currency/protos"
)

type Currency struct {
	protos.UnimplementedCurrencyServer
	l *log.Logger
}

func NewCurrency(l *log.Logger) *Currency {
	return &Currency{l: l}
}

func (c *Currency) GetRate(ctx context.Context, r *protos.RateRequest) (*protos.RateResponse, error) {
	c.l.Printf("Handle GetRate\nbase = %v \ndestination = %v\n", r.Base, r.Destination)
	return &protos.RateResponse{Rate: 0.5}, nil
}
