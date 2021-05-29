package database

import (
	"context"
	"fmt"
	utilsApi "github.com/incafox/golang-api/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var conn *pgxpool.Pool

func InitDB(testenv bool) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env")
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
	if testenv {
		db_host = "localhost"
	}
	db_url := "postgres://" + db_user + ":" + db_password + "@" + db_host + ":" + db_port + "/" + db_name
	conn, err = pgxpool.Connect(context.Background(), db_url)
	if err != nil {
		log.Fatal("Problem initializing database", db_url)
	}
	/*
		if false {
			_, err = conn.Exec(context.Background(), "truncate Drivers")
			if err != nil {
				log.Println("Failed to truncate table ", err)
			}
		}*/
}

func GetDBConn() *pgxpool.Pool {
	return conn
}

func Seed() {
	createSql := `
    create table Drivers(
      id SERIAL PRIMARY KEY,
			email TEXT UNIQUE,
			username TEXT UNIQUE,
			password TEXT)
      ;
  `
	_, err := conn.Exec(context.Background(), createSql)
	if err != nil {
		fmt.Println("Failed to create table: ", err)
	} else {
		fmt.Println("Tables Created")
	}
}
