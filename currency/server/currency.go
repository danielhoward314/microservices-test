package server

import (
	"context"
	"log"

	"github.com/danielhoward314/microservices-test/currency/data"
	protos "github.com/danielhoward314/microservices-test/currency/protos"
)

type Currency struct {
	protos.UnimplementedCurrencyServer
	l  *log.Logger
	er *data.ExchangeRates
}

func NewCurrency(l *log.Logger, er *data.ExchangeRates) *Currency {
	return &Currency{l: l, er: &data.ExchangeRates{}}
}

func (c *Currency) GetRate(ctx context.Context, r *protos.RateRequest) (*protos.RateResponse, error) {
	b := r.Base.String()
	d := r.Destination.String()
	rate, err := c.er.GetRate(b, d)
	if err != nil {
		return nil, err
	}
	return &protos.RateResponse{Rate: float32(rate)}, nil
}
