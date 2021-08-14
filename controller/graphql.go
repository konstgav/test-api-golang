package controller

import (
	"net/http"
	"test-api-golang/interfaces"
)

type GraphqlController struct {
	interfaces.GraphqlControllerInterface
}

func NewGraphqlController() *GraphqlController {
	return &GraphqlController{}
}

func (c *GraphqlController) Graphql(w http.ResponseWriter, r *http.Request) {

}
