package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
)

type ExchangeRates struct {
	rates map[string]float64
}

func NewRates() (*ExchangeRates, error) {
	er := &ExchangeRates{rates: map[string]float64{}}

	err := er.getRates()

	return er, err
}

func (e *ExchangeRates) GetRate(base, dest string) (float64, error) {
	fmt.Printf("e.rates %v", e.rates)
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", base)
	}

	dr, ok := e.rates[dest]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", dest)
	}
	fmt.Printf("base rate = %v\n destination rate = %v\n", br, dr)
	return dr / br, nil
}

func (e *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected error code 200 got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	md := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&md)

	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}

		e.rates[c.Currency] = r
	}

	e.rates["EUR"] = 1

	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
