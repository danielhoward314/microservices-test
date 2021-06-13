package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// Product defines the structure for an API product
type Product struct {
	ID                  int     `json:"id"`
	Name                string  `json:"name" validate:"required"`
	Description         string  `json:"description"`
	BaseCurrency        string  `json:"base_currency"`
	DestinationCurrency string  `json:"destination_currency"`
	BasePrice           float32 `json:"base_price" validate:"gt=0"`
	DestinationPrice    float32 `json:"destination_price" validate:"gt=0"`
	SKU                 string  `json:"sku" validate:"required,sku"`
	CreatedOn           string  `json:"-"`
	UpdatedOn           string  `json:"-"`
	DeletedOn           string  `json:"-"`
}

func (p *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// Products is a collection of Product
type Products []*Product

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	// if second arg int is >0, then that is a cap on # of matches allowed
	matches := re.FindAllString(fl.Field().String(), -1)
	// there should be one and only one match to the valid sku regex
	return len(matches) == 1
}

func (p *Product) Validate() error {
	v := validator.New()
	// first arg corresponds to what is used in struct tag
	v.RegisterValidation("sku", validateSKU)
	return v.Struct(p)
}

// GetProducts returns a list of products
func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	id := getNextId()
	p.ID = id
	productList = append(productList, p)
}

func getNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func UpdateProduct(id int, p *Product) error {
	_, i, err := findProductById(id)
	if err != nil {
		return err
	}
	p.ID = id
	productList[i] = p
	return nil
}

var ErrProductNotFound = fmt.Errorf("product not found by id")

func findProductById(id int) (*Product, int, error) {
	for i, el := range productList {
		if id == el.ID {
			return el, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

// productList is a hard coded list of products for this
// example data source
var productList = []*Product{
	{
		ID:           0,
		Name:         "Latte",
		Description:  "Frothy milky coffee",
		BasePrice:    2.45,
		BaseCurrency: "USD",
		SKU:          "abc-dfjid-vikng",
		CreatedOn:    time.Now().UTC().String(),
		UpdatedOn:    time.Now().UTC().String(),
	},
	{
		ID:           1,
		Name:         "Espresso",
		Description:  "Short and strong coffee without milk",
		BasePrice:    1.99,
		BaseCurrency: "USD",
		SKU:          "moo-prang-tonka",
		CreatedOn:    time.Now().UTC().String(),
		UpdatedOn:    time.Now().UTC().String(),
	},
}
