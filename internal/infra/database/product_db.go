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

func (p *Product) FindAll(page, limit int, sort string) ([]entity.Product, error) {
	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	offset := (page - 1) * limit

	var (
		query string
		rows  *sql.Rows
		err   error
	)

	if limit != 0 && page != 0 {
		query = "select * from products order by created_at " + sort + " limit ? offset ?"
		stmt, err := p.DB.Prepare(query)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		rows, err = stmt.Query(limit, offset)
		if err != nil {
			return nil, err
		}
	} else {

		query = "select * from products"
		rows, err = p.DB.Query(query)
	}

	defer rows.Close()

	var products []entity.Product

	for rows.Next() {
		var p entity.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, err

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

	stmt, err := p.DB.Prepare("update products set name = ?, price = ? where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Price, product.ID)
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
