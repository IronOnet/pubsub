package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/irononet/publisher/utils"
	"github.com/streadway/amqp"
)



func main(){
	// Load data from a csv file 
	// publish that data to a rabbitmq channel 
	records, err := utils.LoadUserRecords("data/users.csv") 
	if err != nil{
		panic(err)
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") 
	if err != nil{
		log.Fatal("could not open channel: ", err)
	}

	// Create a Channel 
	ch, err := conn.Channel() 
	if err != nil{
		log.Fatal("could not create a new channel", err)
	}
	defer ch.Close() 

	q, err := ch.QueueDeclare(
		"users", // queue name 
		false, // Durable 
		false, // Delete when unused
		false, // Exclusive
		false, // no-wait 
		nil, // Arguments 
	)

	if err != nil{
		log.Fatal("error could not start the queue: ", err) 
	}

	for{
		for _, record := range records{
			userJSON, err := json.Marshal(record)
			if err != nil{
				log.Println("Error marshalling user: ", err) 
				continue 
			}
	
			// Publish the serialized data to RabbitMQ 
			err = ch.Publish(
				"", 
				q.Name, // Routing key (queue name)  
				false, // Mandatory 
				false, // Immediate 
				amqp.Publishing{
					ContentType: "application/json", 
					Body: userJSON,
				})
	
			if err != nil{
				log.Println("Error publishing user: ", err) 
			} else{
				fmt.Println("user published: ", record.ID)
			}
		}
	}
	
}