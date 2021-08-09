package main

type CrudServiceInterface interface {
	GetItem(id int) (EntityInterface, error)
	GetList(parameters ListParametersInterface) (EntityInterface, error)
	Create(item EntityInterface) (EntityInterface, error)
	Update(item EntityInterface) (EntityInterface, error)
	Delete(id int) error
}

type CrudService struct {
	CrudServiceInterface
	repository CrudRepositoryInterface
}

func NewCrudService() *CrudService {
	return &CrudService{
		CrudServiceInterface: nil,
		repository:           NewCrudRepository(),
	}
}

func (c CrudService) GetItem(id int) (EntityInterface, error) {
	return c.repository.Find(id)
}

func (c CrudService) GetList(parameters ListParametersInterface) (EntityInterface, error) {
	return c.repository.List(parameters)
}

func (c CrudService) Create(item EntityInterface) (EntityInterface, error) {
	return c.repository.Create(item)
}

func (c CrudService) Update(item EntityInterface) (EntityInterface, error) {
	return c.repository.Update(item)
}

func (c CrudService) Delete(id int) error {
	return c.repository.Delete(id)
}
