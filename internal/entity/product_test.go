package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("Notebook", 100.0)

	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, 100.0, product.Price)
	assert.Equal(t, "Notebook", product.Name)
}

func TestProductWhenNameIsRequired(t *testing.T) {
	product, err := NewProduct("", 100.0)

	assert.Nil(t, product)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestProductWhenPriceIsRequired(t *testing.T) {
	product, err := NewProduct("Notebook", 0.0)

	assert.Nil(t, product)
	assert.Equal(t, ErrPriceIsRequired, err)

}

func TestProductWhenPriceIsInvalid(t *testing.T) {
	product, err := NewProduct("Notebook", -100.0)

	assert.Nil(t, product)
	assert.Equal(t, ErrInvalidPrice, err)
}

func TestProductValidate(t *testing.T) {
	product, err := NewProduct("Notebook", 100.0)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Nil(t, product.Validate())

}
