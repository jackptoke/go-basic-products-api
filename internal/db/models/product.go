package models

import (
	"database/sql"
)

type Product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func GetProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		return []Product{}, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return []Product{}, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (p *Product) GetProduct(db *sql.DB, id int) error {
	row := db.QueryRow("SELECT * FROM products WHERE id=?", id)
	err := row.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) CreateProduct(db *sql.DB) error {
	result, err := db.Exec("INSERT INTO products(name, quantity, price) VALUES(?, ?, ?)", p.Name, p.Quantity, p.Price)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.Id = int(id)
	return nil
}

func (p *Product) UpdateProduct(db *sql.DB, id int) error {
	result, err := db.Exec("UPDATE products set name = ?, quantity = ?, price = ? where id = ?", p.Name, p.Quantity, p.Price, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (p *Product) DeleteProduct(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
