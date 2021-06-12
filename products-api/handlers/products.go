package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/danielhoward314/microservices-test/products-api/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to encode products into JSON", http.StatusInternalServerError)
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
