package network

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize setup the application and connect the database
func (a *App) Initialize() error {
	connectionString := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", DbUser, DbPass, DbHost, DbName)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	a.Router = mux.NewRouter().StrictSlash(true)
	a.handleRoutes()
	return nil
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func sendResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response, err := json.Marshal(data)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(response)
	if err != nil {
		return
	}
}

// sendError will send error response back to the client when the request is not successful
func sendError(w http.ResponseWriter, statusCode int, err error) {
	errMessage := map[string]string{"error": err.Error()}
	sendResponse(w, statusCode, errMessage)
}

func (a *App) handleRoutes() {
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/products/{id}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/products", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/products/{id}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/products/{id}", a.deleteProduct).Methods("DELETE")
}
