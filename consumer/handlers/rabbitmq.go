package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/irononet/consumer/store"
	"github.com/streadway/amqp"
)

type UserMessage struct {
	ID           int   `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	CreatedAt    string `json:"created_at"`
	DeletedAt    string `json:"deleted_at"`
	MergedAt     string `json:"merged_at"`
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

		createdAt, err := convertStringToUnixTimeStamp(userMsg.CreatedAt) 
		log.Println(createdAt)
		if err != nil{
			panic(err) 
		}
		// deletedAt, err := convertStringToUnixTimeStamp(userMsg.DeletedAt) 
		// if err != nil{
		// 	panic(err) 
		// }
		// mergedAt, err := convertStringToUnixTimeStamp(userMsg.MergedAt) 
		// if err != nil{
		// 	panic(err) 
		// }

		// parentUserId, err := strconv.ParseInt(userMsg.ParentUserId, 10, 64) 
		// if err != nil{
		// 	panic(err)
		// }

		user := &store.UserSql{
			ID: uint(userMsg.ID), 
			FirstName: userMsg.FirstName, 
			LastName: userMsg.LastName, 
			EmailAddress: userMsg.EmailAddress, 
			//CreatedAt: createdAt, 
			DeletedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC), 
			MergedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),  
			ParentUserId: 8,
		}

		err = userStore.CreateUser(context.Background(), user)
		if err != nil{
			log.Println("failed to save user:", err) 
			return errors.New(fmt.Sprintf("failed to save user: %v", err))
		}
	}
	return nil 
	
}

func convertStringToUnixTimeStamp(timeStr string) (time.Time, error){
	seconds, err := strconv.ParseInt(timeStr, 10, 64) 
	if err != nil{
		return time.Time{}, err 
	}

	t := time.Unix(seconds, 0) 
	return t, nil 
}

