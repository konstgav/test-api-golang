package service

import (
	"test-api-golang/interfaces"
)

type CrudService struct {
	interfaces.CrudServiceInterface
	repository interfaces.CrudRepositoryInterface
}

func NewCrudService(repository interfaces.CrudRepositoryInterface) *CrudService {
	return &CrudService{
		repository: repository,
	}
}

func (c CrudService) GetItem(id int) (interfaces.EntityInterface, error) {
	return c.repository.Find(id)
}

func (c CrudService) GetList(parameters interfaces.ListParametersInterface) (interfaces.EntityInterface, error) {
	return c.repository.List(parameters)
}

func (c CrudService) Create(item interfaces.EntityInterface) (interfaces.EntityInterface, error) {
	return c.repository.Create(item)
}

func (c CrudService) Update(item interfaces.EntityInterface) (interfaces.EntityInterface, error) {
	return c.repository.Update(item)
}

func (c CrudService) Delete(id int) error {
	return c.repository.Delete(id)
}
