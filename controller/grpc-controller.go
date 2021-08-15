package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"test-api-golang/interfaces"
	"test-api-golang/mailserver"
	"test-api-golang/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type GrpcClientController struct {
	interfaces.GrpcClientControllerInterface
	collection *mongo.Collection
}

func NewGrpcClientController(collection *mongo.Collection) *GrpcClientController {
	return &GrpcClientController{
		collection: collection,
	}
}

type SendMailParams struct {
	ID    int    `json:"_id,omitempty"`
	Email string `json:"email,omitempty"`
}

func (c *GrpcClientController) GetNameById(id int) (string, error) {
	var product model.Product
	filter := bson.M{"_id": id}
	err := c.collection.FindOne(context.TODO(), filter).Decode(&product)
	return product.Name, err
}

func (c *GrpcClientController) SendMail(w http.ResponseWriter, r *http.Request) {
	var params SendMailParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		GetError(err, w)
		return
	}
	name, err := c.GetNameById(params.ID)
	if err != nil {
		GetError(err, w)
		return
	}
	conn, err := grpc.Dial("localhost:20100", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := mailserver.NewMailerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rply, err := client.SendPass(ctx, &mailserver.MsgRequest{
		Email:       params.Email,
		ProductName: name,
	})
	if err != nil {
		log.Println("something went wrong", err)
	}
	log.Println("The mail was send: ", rply.IsSent)
}
