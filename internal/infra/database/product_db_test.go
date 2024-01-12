package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"

	"github.com/pedrohenmonteiro/golang-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	db, err := openDBAndCreateProductTable()
	assert.NoError(t, err)
	defer dropTableProductAndCloseDB(db)

	productDB := NewProduct(db)
	p, err := entity.NewProduct("Notebook", 100.0)
	assert.NoError(t, err)

	err = productDB.Create(p)
	assert.Nil(t, err)
	assert.NotEmpty(t, p.ID)

}

func TestFindAllProducts(t *testing.T) {
	db, err := openDBAndCreateProductTable()
	assert.NoError(t, err)
	defer dropTableProductAndCloseDB(db)

	productDB := NewProduct(db)

	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.NoError(t, err)
		err = productDB.Create(product)
		assert.NoError(t, err)
	}

	products, err := productDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = productDB.FindAll(3, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 23", products[2].Name)
}

func TestFindProductByID(t *testing.T) {
	db, err := openDBAndCreateProductTable()
	assert.NoError(t, err)
	defer dropTableProductAndCloseDB(db)

	product, err := entity.NewProduct("Product 1", 100.0)
	assert.NoError(t, err)

	productDB := NewProduct(db)
	productDB.Create(product)

	product, err = productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Product 1", product.Name)

}

func TestUpdateProduct(t *testing.T) {
	db, err := openDBAndCreateProductTable()
	assert.NoError(t, err)
	defer dropTableProductAndCloseDB(db)

	product, err := entity.NewProduct("Product 1", 100.0)
	assert.NoError(t, err)

	productDB := NewProduct(db)
	productDB.Create(product)

	product.Name = "Product Changed"
	product.Price = 1.00

	err = productDB.Update(product)
	assert.NoError(t, err)

	productFound, err := productDB.FindByID(product.ID.String())
	assert.NoError(t, err)

	assert.Equal(t, "Product Changed", productFound.Name)
	assert.Equal(t, 1.00, productFound.Price)

}

func TestDeleteProduct(t *testing.T) {
	db, err := openDBAndCreateProductTable()
	assert.NoError(t, err)
	defer dropTableProductAndCloseDB(db)

	product, err := entity.NewProduct("Product 1", 10.0)
	assert.NoError(t, err)

	productDB := NewProduct(db)
	productDB.Create(product)

	err = productDB.Delete(product.ID.String())
	assert.NoError(t, err)

	product, err = productDB.FindByID(product.ID.String())
	assert.Error(t, err)
	assert.Nil(t, product)

}

//auxiliares

func openDBAndCreateProductTable() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE products (
			id TEXT PRIMARY KEY,
			name TEXT,
			price REAL,
			created_at DATE
		)
	`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func dropTableProductAndCloseDB(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS products")
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}
