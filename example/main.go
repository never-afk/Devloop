package main

import (
	"github.com/never-afk/Devloop/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

func main() {
	client, err := mongodb.Connect(&mongodb.Config{
		TimeOut:      time.Second * 3,
		DatabaseName: "Devloop",
		DSN:          `mongodb://127.0.0.1`,
	})
	if err != nil {
		print(err)
	}
	table := client.Table(`test`)

	// insert
	table.InsertOne(bson.M{"x": 1})

	docs := []interface{}{
		bson.M{"x": 2},
		bson.M{"x": 3},
	}
	// InsertMany
	table.InsertMany(docs)
	printData(table)

	// UpdateOne
	table.UpdateOne(bson.M{"x": 1}, bson.M{"$set": bson.M{"x": 99}})
	printData(table)

	table.UpdateMany(bson.M{"x": bson.M{"$ne": 99}}, bson.M{"$set": bson.M{"x": 6}})
	printData(table)

	// DeleteMany
	table.DeleteMany(bson.M{})
}

func printData(table *mongodb.CollectionHandel) {
	rs, _ := table.Find(bson.M{})
	result := make([]bson.M, 0)
	rs.All(nil, &result)
	log.Println(result)
}
