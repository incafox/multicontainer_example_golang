package ModelDriver

import (
	"context"
	"github.com/incafox/golang-api/database"
	utilsApi "github.com/incafox/golang-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type User struct {
	Email    string `json:"email" `
	Username string `json:"username" `
	Password string `json:"password"`
}

type RequestRegister struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type RequestLogin struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type RequestToken struct {
	Token string `json:"token"`
}

func GetAllDriversFromMongo(page int, limit int) (interface{}, bool) {
	skipOpt := int64(page)
	limitOpt := int64(limit)
	opts := options.FindOptions{
		Skip:  &skipOpt,
		Limit: &limitOpt,
	}
	cursor, err := database.UseCollection("drivers").Find(context.Background(), bson.D{}, &opts)
	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		log.Println(utilsApi.TypeLog().Error, err)
	}
	if err != nil {
		log.Println(utilsApi.TypeLog().Error, err)
		return results, false
	}
	return results, true
}

func SetDriverAvailability(username string, available int) {
	filter := bson.D{{"username", username}}
	update := bson.D{
		{"$set", bson.D{
			{"available", available},
		}},
	}
	_, err := database.UseCollection("drivers").UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(utilsApi.TypeLog().DatabaseError, err)
	} else {
		log.Println(utilsApi.TypeLog().DatabaseInfo, "driver", username, "updated")
	}
}

func SetDriverDistance(username string, distanceKm float64) {
	filter := bson.D{{"username", username}}
	update := bson.D{
		{"$set", bson.D{
			{"distance", distanceKm},
		}},
	}
	_, err := database.UseCollection("drivers").UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println(utilsApi.TypeLog().DatabaseError, err)
	} else {
		log.Println(utilsApi.TypeLog().DatabaseInfo, "Driver", username, "updated")
	}
}

func GetAllDriversByDistanceFromMongo(distanceKm float64, available int) (interface{}, bool) {
	cursor, err := database.UseCollection("drivers").Find(context.Background(), bson.D{{"distance", bson.D{{"$lt", distanceKm}}}})
	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		log.Println(utilsApi.TypeLog().Error, err)
	}
	if err != nil {
		log.Println(utilsApi.TypeLog().Error, err)
		return results, false
	}
	return results, true
}
