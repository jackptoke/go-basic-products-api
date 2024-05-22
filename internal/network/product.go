package network

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"inventory/internal/db/models"
	"log"
	"net/http"
	"strconv"
)

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := models.GetProducts(a.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err)
	}
	sendResponse(w, http.StatusOK, products)
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	var product models.Product

	err = product.GetProduct(a.DB, productID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			sendError(w, http.StatusNotFound, errors.New("product not found"))
		default:
			sendError(w, http.StatusInternalServerError, err)
		}
		return
	}

	sendResponse(w, http.StatusOK, product)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		sendError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	err := product.CreateProduct(a.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err)
	}
	sendResponse(w, http.StatusCreated, product)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}
	log.Println("Product ID: ", productID)

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		sendError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	product.Id = productID

	err = product.UpdateProduct(a.DB, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, http.StatusNotFound, errors.New("product not found"))
		} else {
			sendError(w, http.StatusInternalServerError, err)
		}
		return
	}

	sendResponse(w, http.StatusOK, product)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, errors.New("invalid product id"))
	}

	var product models.Product

	err = product.GetProduct(a.DB, productID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			sendError(w, http.StatusNotFound, errors.New("product not found"))
		default:
			sendError(w, http.StatusInternalServerError, err)
		}
		return
	}

	err = product.DeleteProduct(a.DB, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, http.StatusNotFound, errors.New("product data has not been deleted"))
		} else {
			sendError(w, http.StatusInternalServerError, err)
		}
		return
	}
	sendResponse(w, http.StatusNoContent, nil)
}
