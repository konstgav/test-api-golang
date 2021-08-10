package main

import (
	"sync"
	"test-api-golang/controller"
	"test-api-golang/interfaces"
	"test-api-golang/repository"
	"test-api-golang/service"
)

type ServiceContainerInterface interface {
	InjectCrudController() interfaces.CrudControllerInterface
}

type kernel struct{}

func (k *kernel) InjectCrudController() interfaces.CrudControllerInterface {
	collection := ConnectDB()
	crudRepository := repository.NewCrudRepository(collection)
	crudService := service.NewCrudService(crudRepository)
	crudController := controller.NewCrudController(crudService)
	return crudController
}

var (
	k             *kernel
	containerOnce sync.Once
)

func ServiceContainer() ServiceContainerInterface {
	if k == nil {
		containerOnce.Do(func() {
			k = &kernel{}
		})
	}
	return k
}
