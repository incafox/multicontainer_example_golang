package ModelDriver

import (
	"context"
	"fmt"
	"github.com/incafox/golang-api/database"
	"log"
	"strings"
)

type DriverDetails struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestLogin struct {
	//Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetDriverByUsername(username string, password string) (DriverDetails, bool) {
	var details DriverDetails
	err := database.GetDBConn().QueryRow(context.Background(), "Select * from Drivers where username = $1 AND password = $2", username, password).Scan(&details.Id, &details.Email, &details.Username, &details.Password)
	fmt.Println(details)
	if err != nil {
		return details, false
	}
	return details, true
}

func CreateDriver(username string, email string, password string) error {
	_, err := database.GetDBConn().Exec(context.Background(), "INSERT INTO Drivers(email, username,password) VALUES ($1,$2,$3)",
		email, username, password)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.Println(err)
			return err
		}
	}
	return nil
}
