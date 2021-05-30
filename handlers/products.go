package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/danielhoward314/microservices-test/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// exposes the endpoint `/products`
// checks the request's method and routes to appropriate handler method
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.GetProducts(rw, r)
		return
	case http.MethodPost:
		p.AddProduct(rw, r)
		return
	case http.MethodPut:
		// pure Go way of deriving path param in URI: use regex
		regex := regexp.MustCompile(`/([0-9]+)`)
		g := regex.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			p.l.Printf("Should only have 1 capture group %v", g)
			http.Error(rw, "Invalid path parameter", http.StatusBadGateway)
			return
		}
		// defensive coding: capture group should have two chars `/` and `<idInt>`
		if len(g[0]) != 2 {
			p.l.Printf("Capture group should have length of 2 for the two chars `/<idInt>` %v", g[0])
			http.Error(rw, "Invalid path parameter", http.StatusNotFound)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Printf("Could not convert string id path param to int %v", g[0][1])
			http.Error(rw, "Invalid path parameter", http.StatusBadRequest)
			return
		}
		p.UpdateProduct(id, rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to encode products into JSON", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle POST /products")
	newProduct := &data.Product{}
	err := newProduct.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to decode new product JSON from request body", http.StatusNotFound)
	}
	p.l.Printf("Product: %#v", newProduct)
	data.AddProduct(newProduct)
}

func (p Products) UpdateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle PUT /products/:id")
	update := &data.Product{}
	err := update.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to decode product update JSON from request body", http.StatusNotFound)
	}
	p.l.Printf("Product: %#v", update)
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
