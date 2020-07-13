package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func GetDBcollection()(*mongo.Collection,*mongo.Collection,*mongo.Collection, error){
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://testing:testing123@cosc4353.x1ecy.mongodb.net/<dbname>?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal()
		return nil, nil,nil, err
	} else {
		fmt.Println("Connection established")
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, nil,nil, err
	}
	collection := client.Database("userLogins").Collection("users")
	deliveryCollection := client.Database("deliveries").Collection("userOrderSummaries")
	userInfo := client.Database("personalInfo").Collection("information")
	return collection, deliveryCollection, userInfo, nil
}
