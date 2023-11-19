package database

import (
	"fmt"
	"testing"

	"github.com/reinaldosaraiva/go-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
func setupDatabase(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
    if err != nil {
        t.Fatal(err)
    }
    db.AutoMigrate(&entity.Product{})
    return db
}

func  TestCreateNewProduct(t *testing.T){
	db := setupDatabase(t)
	product,err := entity.NewProduct("Product 1", 10)
	assert.NoError(t, err)
	productDB := NewProduct(db)
	err = productDB.Create(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
}

func TestFindAllProducts(t *testing.T){
	db := setupDatabase(t)
	for i := 1; i < 24; i++{
		product,err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Intn(100)*10)
		assert.NoError(t, err)
		productDB := NewProduct(db)
		err = productDB.Create(product)
		assert.NoError(t, err)
	}
	productDB := NewProduct(db)
	products,err := productDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, products[0].Name, "Product 1")
	assert.Equal(t, products[9].Name, "Product 10")
}

func TestFindProductByID(t *testing.T){
	db := setupDatabase(t)
	product,err := entity.NewProduct("Product 1", 10)
	assert.NoError(t, err)
	db.Create(product)
	productDB := NewProduct(db)
	p,err := productDB.FindByID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, p.Name, "Product 1")
}

func TestUpdateProduct(t *testing.T){
	db := setupDatabase(t)
	product,err := entity.NewProduct("Product 1", 10)
	assert.NoError(t, err)
	db.Create(product)
	productDB := NewProduct(db)
	product.Name = "Product 2"
	err = productDB.Update(product)
	assert.NoError(t, err)
	p,err := productDB.FindByID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, p.Name, "Product 2")
}

func TestDeleteProduct(t *testing.T){
	db := setupDatabase(t)
	product,err := entity.NewProduct("Product 1", 10)
	assert.NoError(t, err)
	db.Create(product)
	productDB := NewProduct(db)
	err = productDB.Delete(product.ID)
	assert.NoError(t, err)
	p2,err := productDB.FindByID(product.ID)
	assert.Error(t, err)
	assert.Nil(t, p2)
}
