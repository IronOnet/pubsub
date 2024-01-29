package handlers

import (
	//"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Event keeps a list of clients
// who are currently attached and
// broadcasts events to those clients
type Event struct{
	Message chan string 
	NewClients chan chan string 
	ClosedClients chan chan string 

	TotalClients map[chan string]bool 
}

type ClientChan chan string 

func NewServer() (*Event){
	event := &Event{
		Message: make(chan string), 
		NewClients: make(chan chan string), 
		ClosedClients: make(chan chan string), 
		TotalClients: make(map[chan string]bool), 
	}

	go event.listen() 

	return event 
}

func (stream Event) listen(){
	for{
		select{
		case client := <-stream.NewClients: 
			stream.TotalClients[client] = true 
			log.Printf("client added. %d registered clients", len(stream.TotalClients))

		// Remove closed client 
		case client := <-stream.ClosedClients: 
			delete(stream.TotalClients, client) 
			close(client) 
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		case eventMsg := <-stream.Message: 
			for clientMessageChan := range stream.TotalClients{
				clientMessageChan <- eventMsg 
			}
		}

		
	}
}

func (stream *Event) ServeHTTP() gin.HandlerFunc{
	return func(c *gin.Context){
		clientChan := make(ClientChan)

		stream.NewClients <- clientChan 

		defer func(){
			stream.ClosedClients <- clientChan 
		}()

		c.Set("clientChan", clientChan)
		c.Next()
	}
}

func HeaderMiddleWare() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Header().Set("Content-Type", "text/event-stream") 
		c.Writer.Header().Set("Cache-Control", "no-cache") 
		c.Writer.Header().Set("Connection", "keep-alive") 
		c.Writer.Header().Set("Transfer-Enconding", "chunked")
		c.Next()
	}
}

func StreamEvents(c *gin.Context){
	// Query the API and send the users 
	// As SSE in realtime 
	c.Writer.Header().Set("Content-Type", "text/event-stream") 

	for i := 0; i < 1000; i++{
		//c.SSEvent("message", fmt.Sprintf("stream route %d", i))
		c.JSON(http.StatusOK, gin.H{"stream": i}) 
		time.Sleep(1 * time.Second)
	}
}