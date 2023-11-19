package entity

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNameRequired = errors.New("Name is required")
	ErrPriceInvalid = errors.New("Price is invalid")
)
type Product struct {
	Name string `json:"name"`
	Price int `json:"price"`
	gorm.Model
}

func NewProduct(name string, price int) (*Product,error) {
	product := &Product{
		Name: name,
		Price: price,
	}
	err := product.Validate()
	if err != nil {
		return nil,err
	}
	return product,nil
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrNameRequired
	}
	if p.Price <= 0 {
		return ErrPriceInvalid
	}
	return nil
}