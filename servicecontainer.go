package main

import (
	"os"
	"sync"
	"test-api-golang/controller"
	"test-api-golang/interfaces"
	"test-api-golang/redis"
	"test-api-golang/repository"
	"test-api-golang/service"
	"time"
)

const (
	RedisRecordTimeToLiveSeconds = 5
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
	redis_URI := os.Getenv("REDIS_URI")
	if redis_URI == "" {
		panic("Environmental variable REDIS_URI do not set")
	}
	redisCache := redis.NewCache(redis_URI, "", 0, time.Duration(RedisRecordTimeToLiveSeconds)*time.Second)
	crudController := controller.NewCrudController(crudService, redisCache)
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
