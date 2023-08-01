package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sano/config"
	"time"
)

type Lookup struct {
	Name        *string   `json:"name,omitempty"`
	DisplayName *string   `json:"display_name"`
	Status      bool      `json:"status"`
	Time        time.Time `json:"time"`
}

type ServiceLookup struct {
	Name    string   `json:"name"`
	Lookups []Lookup `json:"lookups"`
}

var Client *mongo.Client

func ConnectToDatabase(ctx context.Context, url string) *mongo.Client {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPIOptions)

	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return Client
}

func StoreLookup(service config.Service, online bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := Client.Database("test").Collection("services")
	_, err := collection.InsertOne(ctx, bson.M{"name": service.Name, "online": online, "time": time.Now()})
	if err != nil {
		log.Panicln(err)
	}
}

func GetLookups() []Lookup {
	var lookups []Lookup

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := Client.Database("test").Collection("services")
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", bson.D{{"service", "$name"}}},
			{"online", bson.D{{"$push", "$online"}}},
			{"time", bson.D{{"$push", "$time"}}},
		}}},
		{{"$sort", bson.D{
			{"time", -1},
		}}},
	}

	cur, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Panicln(err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		_id := result["_id"].(bson.M)
		service := _id["service"].(string)
		online := result["online"].(bson.A)
		onlineTime := result["time"].(bson.A)

		var displayName *string
		configService := config.Config.GetService(service)
		if configService == nil || (configService != nil && configService.DisplayName == nil) {
			displayName = nil
		} else {
			displayName = configService.DisplayName
		}

		lookups = append(lookups, Lookup{
			Name:        &service,
			DisplayName: displayName,
			Status:      online[len(online)-1].(bool),
			Time:        onlineTime[(len(onlineTime) - 1)].(primitive.DateTime).Time(),
		})
	}

	log.Println(lookups[1])

	if err := cur.Err(); err != nil {
		log.Panicln(err)
	}

	return lookups
}

func GetLookup(serviceName string) ServiceLookup {
	var serviceLookup ServiceLookup

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := Client.Database("test").Collection("services")
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"name", serviceName}}}},
		{{"$group", bson.D{
			{"_id", bson.D{{"service", "$name"}}},
			{"online", bson.D{{"$push", "$online"}}},
			{"time", bson.D{{"$push", "$time"}}},
		}}},
		{{"$sort", bson.D{
			{"time", -1},
		}}},
	}
	cur, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Panicln(err)
	}
	defer cur.Close(ctx)
	cur.Next(ctx)

	var result bson.M
	err = cur.Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	_id := result["_id"].(bson.M)
	online := result["online"].(bson.A)
	onlineTime := result["time"].(bson.A)

	var lookups []Lookup
	for i := 0; i < len(online); i++ {
		lookups = append(lookups, Lookup{
			Status: online[i].(bool),
			Time:   onlineTime[i].(primitive.DateTime).Time(),
		})
	}

	serviceLookup = ServiceLookup{
		Name:    _id["service"].(string),
		Lookups: lookups,
	}

	if err := cur.Err(); err != nil {
		log.Panicln(err)
	}

	return serviceLookup
}
