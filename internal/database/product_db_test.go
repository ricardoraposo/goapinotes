package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ricardoraposo/api-again/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateNewProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})
	product, err := entity.NewProduct("product", 100)
	assert.Nil(t, err)
	productDB := NewProduct(db)
	err = productDB.Create(product)
	assert.Nil(t, err)
	assert.NotEmpty(t, product.ID)
}

func TestFindAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})
	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("product %d", i), rand.Intn(100))
        assert.NoError(t, err)
        db.Create(product)
	}

    productDB := NewProduct(db)
    products, err := productDB.FindAll(1, 10, "asc")
    assert.NoError(t, err)
    assert.Len(t, products, 10)
    assert.Equal(t, "product 1", products[0].Name)
    assert.Equal(t, "product 10", products[9].Name)

    products, err = productDB.FindAll(2, 10, "asc")
    assert.NoError(t, err)
    assert.Len(t, products, 10)
    assert.Equal(t, "product 11", products[0].Name)
    assert.Equal(t, "product 20", products[9].Name)

    products, err = productDB.FindAll(3, 10, "asc")
    assert.NoError(t, err)
    assert.Len(t, products, 3)
    assert.Equal(t, "product 21", products[0].Name)
    assert.Equal(t, "product 23", products[2].Name)
}
