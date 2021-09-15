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
		err = c.cache.Set(strconv.Itoa(id), product)
		if err != nil {
			log.Println(err.Error())
		}
	}
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		GetError(err, w)
	}
}

func (c CrudController) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var products, err = c.service.GetList(nil)
	if err != nil {
		GetError(err, w)
		return
	}
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		GetError(err, w)
	}
}

func (c CrudController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var entity interfaces.EntityInterface
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		GetError(err, w)
		return
	}
	result, err := c.service.Create(entity)
	if err != nil {
		GetError(err, w)
		return
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		GetError(err, w)
	}
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
	err = json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		GetError(err, w)
		return
	}

	result, err := c.service.Update(id, entity)
	if err != nil {
		GetError(err, w)
		return
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		GetError(err, w)
	}
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
	err = json.NewEncoder(w).Encode("204 No Data")
	if err != nil {
		GetError(err, w)
	}
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
	_, err = w.Write(message)
	if err != nil {
		log.Println(err.Error())
	}
}
