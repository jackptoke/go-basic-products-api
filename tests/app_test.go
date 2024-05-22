package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"inventory/internal/network"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var app network.App

func TestMain(m *testing.M) {
	err := app.Initialize(network.DbUser, network.DbPass, network.DbHost, "testDb")
	if err != nil {
		log.Fatal(err)
	}

	createTable()

	m.Run()
}

func createTable() {
	query := `CREATE TABLE IF NOT EXISTS products (
    id int NOT NULL AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    quantity int NOT NULL,
    price int NOT NULL,
    PRIMARY KEY (id)
);`
	_, err := app.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	app.DB.Exec("DELETE FROM products")
}

func addProduct(name string, quantity int, price float64) int {
	q := fmt.Sprintf(`INSERT INTO products (name, quantity, price) values ('%v', %v, %v)`, name, quantity, price)
	result, err := app.DB.Exec(q)
	if err != nil {
		log.Fatal(err)
	}
	count, _ := result.RowsAffected()
	log.Printf("Rows affected: %d", count)
	id, err := result.LastInsertId()
	if err != nil {
		return 0
	}
	return int(id)
}

func TestGetProducts(t *testing.T) {
	clearTable()
	id := addProduct("Macbook 16", 5, 5600)
	r, _ := http.NewRequest("GET", fmt.Sprintf("/products/%v", id), nil)
	response := sendRequest(r)
	log.Printf("Response: %v", response.Result().StatusCode)
	checkStatusCode(t, http.StatusOK, response.Result().StatusCode)
}

// sendRequest is a helper function sending request to the test server instance
func sendRequest(r *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	app.Router.ServeHTTP(recorder, r)
	return recorder
}

// checkStatusCode helps check if the test is successful
func checkStatusCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestPostProducts(t *testing.T) {
	clearTable()
	product := []byte(`{"name": "iPhone", "quantity":5, "price": 850}`)

	r, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(product))
	r.Header.Set("Content-Type", "application/json")
	response := sendRequest(r)
	log.Printf("Response: %v", response.Result().StatusCode)
	checkStatusCode(t, http.StatusCreated, response.Result().StatusCode)

	var m map[string]interface{}

	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		return
	}
	if m["name"] != "iPhone" {
		t.Errorf("Expected name to be 'iPhone'. Got '%v'", m["name"])
	}

	if m["quantity"] != 5.0 {
		t.Errorf("Expected quantity to be '5'. Got '%v'", m["quantity"])
	}
	if m["price"] != 850.0 {
		t.Errorf("Expected price to be '850'. Got '%v'", m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	id := addProduct("Macbook 16", 5, 5600)

	r, _ := http.NewRequest("DELETE", fmt.Sprintf("/products/%v", id), nil)
	response := sendRequest(r)
	log.Printf("Response: %v", response.Result().StatusCode)
	checkStatusCode(t, http.StatusNoContent, response.Result().StatusCode)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	id := addProduct("Macbook 16", 5, 5600)

	product := []byte(`{"name": "Macbook 13", "quantity":6, "price": 5995}`)

	r, _ := http.NewRequest("PUT", fmt.Sprintf("/products/%v", id), bytes.NewBuffer(product))
	r.Header.Set("Content-Type", "application/json")
	response := sendRequest(r)
	log.Printf("Response: %v", response.Result().StatusCode)
	checkStatusCode(t, http.StatusOK, response.Result().StatusCode)

	var m map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		log.Fatal(err)
		return
	}
	if m["name"] != "Macbook 13" {
		t.Errorf("Expected name to be 'Macbook 13'. Got '%v'", m["name"])
	}
	if m["quantity"] != 6.0 {
		t.Errorf("Expected quantity to be '6'. Got '%v'", m["quantity"])
	}
	if m["price"] != 5995.0 {
		t.Errorf("Expected price to be '5995'. Got '%v'", m["price"])
	}
}
