package database

import (
	"context"
	utilsApi "github.com/incafox/golang-api/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var database *mongo.Database

func InitDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(utilsApi.TypeLog().Error, "Error loading .env")
	}
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")
	utilsApi.DefaultEnvField(db_user, "root")
	utilsApi.DefaultEnvField(db_password, "root")
	utilsApi.DefaultEnvField(db_host, "localhost")
	utilsApi.DefaultEnvField(db_port, "5432")
	utilsApi.DefaultEnvField(db_name, "test")
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	if os.Getenv("testenv") == "true" {
		db_host = "localhost"
	}
	clientOptions := options.Client().ApplyURI("mongodb://" + db_host + ":" + db_port)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Println(utilsApi.TypeLog().DatabaseError, "Could not connect to database")
	} else {
		log.Println(utilsApi.TypeLog().DatabaseInfo, "Connected to mongodb")
	}
	database = client.Database(db_name)
}
func UseDatabase() *mongo.Database {
	return database
}

func UseCollection(collName string) *mongo.Collection {
	return database.Collection(collName)
}

func Seed() {
}
