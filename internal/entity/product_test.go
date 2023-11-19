package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestNewProduct(t *testing.T) {
	product, err := NewProduct("Feijão", 1000)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	// assert.NotEmpty(t, product.ID)
	assert.NotEmpty(t, product.Name)
	assert.NotEmpty(t, product.Price)
	assert.Equal(t, product.Name, "Feijão")
	assert.Equal(t, product.Price, 1000)
}

func TestProductWhenNameRequired(t *testing.T) {
	product, err := NewProduct("", 10)
	assert.Nil(t, product)
	assert.Equal(t, ErrNameRequired,err)
}
func TestProductWhenPriceRequired(t *testing.T) {
	product, err := NewProduct("Feijão", 0)
	assert.Nil(t, product)
	assert.Equal(t, ErrPriceInvalid,err)
}