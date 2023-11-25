package main  

import (
	"encoding/csv" 
	"encoding/json" 
	"fmt" 
	"log" 
	"os" 

	"github.com/streadway/ampq"
)


type User struct{
	ID int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	CreatedAt string `json:"created_at"`
	DeletedAt string `json:"deleted_at"`
}
func main(){
	// Load data from a csv file 
	// publish that data to a rabbitmq channel 
}