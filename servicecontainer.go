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
	InjectGrpcClientController() interfaces.GrpcClientControllerInterface
}

type kernel struct{}

func (k *kernel) InjectCrudController() interfaces.CrudControllerInterface {
	collection := repository.ConnectDB()
	crudRepository := repository.NewCrudRepository(collection)
	crudService := service.NewCrudService(crudRepository)
	crudController := controller.NewCrudController(crudService)
	return crudController
}

func (k *kernel) InjectGrpcClientController() interfaces.GrpcClientControllerInterface {
	collection := repository.ConnectDB()
	grpcClientController := controller.NewGrpcClientController(collection)
	return grpcClientController
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
