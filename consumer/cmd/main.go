package main

import (
	"context"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/irononet/consumer/handlers"
	"github.com/irononet/consumer/store"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	WEBSERVER_PORT int = 8000 
	RABBITMQ_ADDRESS string = "amqp://guest:guest@localhost:5672/"
)

var (
	ApiHandler handlers.UserHandler
	UserStore *store.UserStoreSql 
	db *gorm.DB 
	Cache *redis.Client 
)


func init(){
	// Init database
	dsn := "root:password@tcp(localhost:3306)/pubsub?charset=utf8&parseTime=True&loc=Local"
	db = getMysqlDb(dsn) 
	UserStore = &store.UserStoreSql{Db: db}
	Cache = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", 
		Password: "", 
		DB: 0, // use default DB
	})

	ApiHandler = *handlers.NewUserHandler(context.Background(), UserStore, Cache)
}


func main(){

	// Initialize data consumption from  
	// the rabbitmq queue 
	go initRabbitMqChannel() 
	// Initialize the rest api server 
	go initApiServer()
}


func initRabbitMqChannel(){
	log.Println("initializing the rabbitmq channel")
	defer log.Println("rabbitmq channel initialization complete")
	err := handlers.InitChannel(RABBITMQ_ADDRESS, *UserStore)
	if err != nil{
		log.Println("error starting the rabbitmq channel: ", err) 
	}
}

func initApiServer(){
	log.Println("starting the API server at port: " ,WEBSERVER_PORT)

	router := gin.Default() 

	router.GET("/user/:id", ApiHandler.GetOneUserById) 
	router.GET("/user/:email", ApiHandler.GetUserByEmail) 
	router.GET("/users/:limit/:skip", ApiHandler.GetUsers) 

	err := router.Run(":" + strconv.Itoa(WEBSERVER_PORT)) 
	if err != nil{
		log.Fatal("failed to start the api server: ", err) 
	}
	
}

func getMysqlDb(connectionStr string) *gorm.DB{
	dsn := connectionStr 
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil{
		log.Println("could not connect to database an error occurred: ", err) 
		return nil 
	}
	return db 
}

