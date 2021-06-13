package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	protos "github.com/danielhoward314/microservices-test/currency/protos"
	"github.com/danielhoward314/microservices-test/products-api/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l  *log.Logger
	cc protos.CurrencyClient
}

func NewProducts(l *log.Logger, cc protos.CurrencyClient) *Products {
	return &Products{l, cc}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	prods := data.GetProducts()
	err := prods.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to encode products into JSON", http.StatusInternalServerError)
		return
	}
}

func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Invalid product id", http.StatusBadRequest)
		return
	}
	prods := data.GetProducts()
	if (id < 0) || (id >= len(prods)) {
		http.Error(rw, "Invalid product id", http.StatusBadRequest)
		return
	}
	prod := prods[id]
	cur := r.URL.Query().Get("currency")
	if cur == "" {
		cur = "USD"
	}
	destination, dExists := protos.Currencies_value[cur]
	base, bExists := protos.Currencies_value[prod.BaseCurrency]
	if dExists && bExists && (destination != base) {
		baseEnum := protos.Currencies(base)
		destinationEnum := protos.Currencies(destination)
		rr, err := p.cc.GetRate(context.Background(), &protos.RateRequest{Base: baseEnum, Destination: destinationEnum})
		if err != nil {
			p.l.Printf("unable to process currency conversion for base %v and destination %v\n", baseEnum, destinationEnum)
			p.l.Print("reverting to default currency of product")
		} else {
			prod.DestinationCurrency = cur
			prod.DestinationPrice = rr.Rate
		}
	}
	err = prod.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to encode product into JSON", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle POST /products")
	newProduct := &data.Product{}
	err := newProduct.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to decode new product JSON from request body", http.StatusNotFound)
		return
	}
	p.l.Printf("Product: %#v", newProduct)
	err = newProduct.Validate()
	if err != nil {
		http.Error(rw, "Invalid request body JSON", http.StatusBadRequest)
		return
	}
	data.AddProduct(newProduct)
}

func (p Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle PUT /products/:id")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to decode path param `id` from request", http.StatusBadRequest)
		return
	}
	update := &data.Product{}
	err = update.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to decode product update JSON from request body", http.StatusNotFound)
		return
	}
	p.l.Printf("Product: %#v", update)
	err = update.Validate()
	if err != nil {
		http.Error(rw, "Invalid request body JSON", http.StatusBadRequest)
		return
	}
	err = data.UpdateProduct(id, update)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Unable to update product", http.StatusInternalServerError)
		return
	}

}
