package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/irononet/consumer/store"
	"github.com/streadway/amqp"
)

type UserMessage struct {
	ID           int   `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	CreatedAt    time.Time `json:"created_at"`
	DeletedAt    time.Time `json:"deleted_at"`
	MergedAt     time.Time `json:"merged_at"`
	ParentUserId int  `json:"parent_user_id"`
}

func InitChannel(connString string, userStore store.UserStoreSql) error {
	conn, err := amqp.Dial(connString) 
	if err != nil{
		return errors.New(fmt.Sprintf("could not connect to broker an error occurred %s", err)) 
	}
	defer conn.Close() 

	ch, err := conn.Channel() 
	if err != nil{
		log.Fatal("failed to open a channel", err) 
	}
	defer ch.Close() 

	// Declare the queue 
	queueName := "users"  
	q, err := ch.QueueDeclare(
		queueName, 
		false, // Durable 
		false, // Deelete when unused 
		false, // Exclusive 
		false, // No-Wait 
		nil, // Arguments
	)

	if err != nil{
		log.Fatal("failed to declare a queue:", err) 
		return errors.New(fmt.Sprintf("failed to declare a queue an error occurred: %s", err))
	}

	msgs, err := ch.Consume(
		q.Name, 
		"", // Consumser 
		true, // Auto-Ack
		false, // Exclusive 
		false, // No-Local 
		false, // No-Wait 
		nil, // Args
	)

	if err != nil{
		log.Fatal("failed to register a consumer", err) 
		return errors.New(fmt.Sprintf("failed to register a consumer %v", err))
	}

	// Create a instance of store  
	for msg := range msgs{
		var userMsg UserMessage 

		// Unmarshall the JSON message into UserMessage struct 
		err := json.Unmarshal(msg.Body, &userMsg) 
		if err != nil{
			log.Println("failed to unmarshal message:", err) 
			continue 
		}

		user := &store.UserSql{
			ID: uint(userMsg.ID), 
			FirstName: userMsg.FirstName, 
			LastName: userMsg.LastName, 
			EmailAddress: userMsg.EmailAddress, 
			CreatedAt: userMsg.CreatedAt, 
			DeletedAt: userMsg.DeletedAt, 
			MergedAt: userMsg.MergedAt, 
			ParentUserId: uint(userMsg.ParentUserId),
		}

		err = userStore.CreateUser(context.Background(), user)
		if err != nil{
			log.Println("failed to save user:", err) 
			return errors.New(fmt.Sprintf("failed to save user: %v", err))
		}
	}
	return nil 
	
}

