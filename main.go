package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	fmt.Println("bismillah")
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://buds_user:meAq0DDWyuLrDZbg@cluster0.rx3zi.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	client.Ping(ctx, nil)
	err = client.Ping(ctx, readpref.Primary())
	collection := client.Database("testing").Collection("numbers")
	res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
	id := res.InsertedID
	fmt.Println("ID: ", id)

	// Query All
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
		fmt.Println("Result: ", result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Query One
	var result struct {
		Value float64
	}
	str := "pi"
	filter := bson.D{{"name", str}}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		// Do something when no record was found
		fmt.Println("record does not exist")
	} else if err != nil {
		log.Fatal(err)
	}
	// Do something with result...
	fmt.Println("One Result: ", result)

	// Delete
	filter = bson.D{{"name", "pi"}}
	collection.DeleteOne(ctx, filter)

	// Delete many
	filter = bson.D{{}}
	collection.DeleteMany(ctx, filter)
}
