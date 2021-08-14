package graphql

import (
	"errors"
	"log"
	"test-api-golang/interfaces"
	"test-api-golang/model"
	"test-api-golang/repository"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var productType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"sku": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"price": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Get (read) single product by id
			   http://localhost:5000/graphql?query={product(id:2){id,name}}
			*/
			"product": &graphql.Field{
				Type:        productType,
				Description: "Get product by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, ok := params.Args["id"].(int)
					if ok {
						return repo.Find(id)
					}
					return nil, nil
				},
			},
			/* Get (read) product list
			   http://localhost:5000/graphql?query={list{id,name,sku,type,price}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(productType),
				Description: "Get product list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return repo.List(nil)
				},
			},
		},
	})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		/* Create new product item
		http://localhost:5000/graphql?query=mutation+_{create(id:10,name:"armor",sku:"v8",type:"item",price:13){id,name,sku,type,price}}
		*/
		"create": &graphql.Field{
			Type:        productType,
			Description: "Create new product",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"sku": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"type": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				product := model.Product{
					ID:    params.Args["id"].(int),
					Name:  params.Args["name"].(string),
					Sku:   params.Args["sku"].(string),
					Type:  params.Args["type"].(string),
					Price: params.Args["price"].(int),
				}
				return repo.Create(product)
			},
		},

		/* Update product by id
		   http://localhost:5000/graphql/product?query=mutation+_{update(id:1,name:"armor",sku:"v8",type:"item",price:13){id,name,sku,type,price}}
		*/
		"update": &graphql.Field{
			Type:        productType,
			Description: "Update product by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"sku": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"type": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, ok := params.Args["id"].(int)
				if !ok {
					return nil, errors.New("need id-field to update data")
				}
				product := make(map[string]interface{})
				product["ID"] = params.Args["id"].(int)
				product["Name"] = params.Args["name"].(string)
				product["Sku"] = params.Args["sku"].(string)
				product["Type"] = params.Args["type"].(string)
				product["Price"] = params.Args["price"].(int)
				return repo.Update(id, product)
			},
		},

		/* Delete product by id
		   http://localhost:5000/graphql/product?query=mutation+_{delete(id:1){id,name,sku,type,price}}
		*/
		"delete": &graphql.Field{
			Type:        productType,
			Description: "Delete product by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, ok := params.Args["id"].(int)
				if ok {
					return "204 No Data", repo.Delete(id)
				}
				return nil, nil
			},
		},
	},
})

func defineSchema() graphql.SchemaConfig {
	return graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	}
}

var repo interfaces.CrudRepositoryInterface

func CreateGrapQLHandler() *handler.Handler {
	collection := repository.ConnectDB()
	repo = repository.NewCrudRepository(collection)
	schema, err := graphql.NewSchema(defineSchema())
	if err != nil {
		log.Panic("Error when creating the graphQL schema", err)
	}
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	return h
}
