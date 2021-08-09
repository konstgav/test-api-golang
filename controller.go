package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CrudControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type CrudController struct {
	CrudControllerInterface
	service CrudServiceInterface
}

func NewCrudController() *CrudController {
	return &CrudController{
		CrudControllerInterface: nil,
		service:                 NewCrudService(),
	}
}

func (c CrudController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		GetError(err, w)
		return
	}
	product, err := c.service.GetItem(id)
	if err != nil {
		GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func (c CrudController) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var products, err = c.service.GetList(nil)
	if err != nil {
		GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(products)
}

func (c CrudController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product Product
	_ = json.NewDecoder(r.Body).Decode(&product)
	result, err := c.service.Create(product)
	if err != nil {
		GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (c CrudController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		GetError(err, w)
		return
	}
	fmt.Println(id)

	var product Product
	_ = json.NewDecoder(r.Body).Decode(&product)
	result, err := c.service.Update(product)

	if err != nil {
		GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (c CrudController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		GetError(err, w)
		return
	}
	err = c.service.Delete(id)
	if err != nil {
		GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode("204 No Data")
}

type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func GetError(err error, w http.ResponseWriter) {
	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
