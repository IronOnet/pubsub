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
	WEBSERVER_PORT   int    = 8000
	RABBITMQ_ADDRESS string = "amqp://guest:guest@rabbitmq/"
)

var (
	ApiHandler handlers.UserHandler
	UserStore  *store.UserStoreSql
	db         *gorm.DB
	Cache      *redis.Client
)

func init() {
	// Init database
	dsn := "mysql:password@tcp(db:3306)/pubsub?charset=utf8&parseTime=True&loc=Local"
	db = getMysqlDb(dsn)
	UserStore = &store.UserStoreSql{Db: db}
	Cache = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0, // use default DB
	})

	ApiHandler = *handlers.NewUserHandler(context.Background(), UserStore, Cache)
}

func main() {

	// Initialize data consumption from
	// the rabbitmq queue
	go initRabbitMqChannel()
	initApiServer()
	// Initialize the rest api server

}

func initRabbitMqChannel() {
	log.Println("initializing the rabbitmq channel")
	defer log.Println("rabbitmq channel initialization complete")
	err := handlers.InitChannel(RABBITMQ_ADDRESS, *UserStore)
	if err != nil {
		log.Println("error starting the rabbitmq channel: ", err)
	}
}

func initApiServer() {
	log.Println("starting the API server at port: ", WEBSERVER_PORT)

	router := gin.Default()

	router.GET("/", ApiHandler.HelloHandler)
	router.GET("/user/:id", ApiHandler.GetOneUserById)
	//router.GET("/user/:email", ApiHandler.GetUserByEmail)
	router.GET("/users/", ApiHandler.GetUsers)

	// SSE stream
	//router.GET("/stream/", handlers.StreamEvents)

	err := router.Run(":" + strconv.Itoa(WEBSERVER_PORT))
	if err != nil {
		log.Fatal("failed to start the api server: ", err)
	}

}

func getMysqlDb(connectionStr string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(connectionStr), &gorm.Config{
		// Optimization tweaks to speed up query write performance
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		log.Println("could not connect to database an error occurred: ", err)
		return nil
	}
	return db
}
