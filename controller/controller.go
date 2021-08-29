package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"test-api-golang/interfaces"

	"github.com/gorilla/mux"
)

type CrudController struct {
	interfaces.CrudControllerInterface
	service interfaces.CrudServiceInterface
	cache   interfaces.Cacher
}

func NewCrudController(service interfaces.CrudServiceInterface, cache interfaces.Cacher) *CrudController {
	return &CrudController{
		service: service,
		cache:   cache,
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
	product, _ := c.cache.Get(strconv.Itoa(id))
	if product == nil {
		product, err = c.service.GetItem(id)
		if err != nil {
			GetError(err, w)
			return
		}
		c.cache.Set(strconv.Itoa(id), product)
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
	var entity interfaces.EntityInterface
	_ = json.NewDecoder(r.Body).Decode(&entity)
	result, err := c.service.Create(entity)
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

	var entity interfaces.EntityInterface
	_ = json.NewDecoder(r.Body).Decode(&entity)
	result, err := c.service.Update(id, entity)

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
	log.Println(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
