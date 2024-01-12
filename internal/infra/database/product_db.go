package database

import (
	"database/sql"

	"github.com/pedrohenmonteiro/golang-api/internal/entity"
)

type Product struct {
	DB *sql.DB
}

func NewProduct(db *sql.DB) *Product {
	return &Product{DB: db}
}

func (p *Product) Create(product *entity.Product) error {
	stmt, err := p.DB.Prepare("insert into products(id, name, price, created_at) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.ID.String(), product.Name, product.Price, product.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *Product) FindByID(id string) (*entity.Product, error) {
	stmt, err := p.DB.Prepare("select * from products where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var product entity.Product

	err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price, &product.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &product, err
}

func (p *Product) Update(product *entity.Product) error {
	_, err := p.FindByID(product.ID.String())
	if err != nil {
		return err
	}

	stmt, err := p.DB.Prepare("update products set name = ?, set price = ? where id = ?")
	if err != nil {
		return err
	}
	stmt.Close()

	_, err = stmt.Exec(product.Name, product.Price)
	if err != nil {
		return err
	}

	return nil
}

func (p *Product) Delete(id string) error {

	_, err := p.FindByID(id)
	if err != nil {
		return err
	}

	stmt, err := p.DB.Prepare("delete from products where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
