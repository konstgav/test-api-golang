package interfaces

import (
	"net/http"
)

type CrudControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type CrudServiceInterface interface {
	GetItem(id int) (EntityInterface, error)
	GetList(parameters ListParametersInterface) (EntityInterface, error)
	Create(item EntityInterface) (EntityInterface, error)
	Update(item EntityInterface) (EntityInterface, error)
	Delete(id int) error
}

type CrudRepositoryInterface interface {
	Find(id int) (EntityInterface, error)
	List(parameters ListParametersInterface) (EntityInterface, error)
	Create(item EntityInterface) (EntityInterface, error)
	Update(item EntityInterface) (EntityInterface, error)
	Delete(id int) error
}

type EntityInterface interface{}

type ListParametersInterface interface{}
