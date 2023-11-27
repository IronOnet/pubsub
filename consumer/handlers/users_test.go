package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/irononet/consumer/store"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	client *redis.Client
)

func TestUserHandler_GetOneUserById(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	mockCache, err := miniredis.Run()
	assert.NoError(t, err)
	cache := redis.NewClient(&redis.Options{
		Addr: mockCache.Addr(),
	})
	sqlStore := &store.UserStoreSql{Db: db}

	err = sqlStore.Migrate(context.Background())
	assert.NoError(t, err)

	user := &store.UserSql{
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
		ParentUserId: 0,
	}

	err = sqlStore.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	handler := NewUserHandler(context.Background(), sqlStore, cache)

	router := gin.New()
	router.GET("/user/:id", handler.GetOneUserById)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["user"])
}

func TestHandler_GetUserByEmail(t *testing.T) {
	// Setup
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	mockCache, err := miniredis.Run()
	assert.NoError(t, err)
	cache := redis.NewClient(&redis.Options{
		Addr: mockCache.Addr(),
	})
	sqlStore := &store.UserStoreSql{Db: db}

	err = sqlStore.Migrate(context.Background())
	assert.NoError(t, err)

	user := &store.UserSql{
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
		ParentUserId: 0,
	}

	err = sqlStore.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	handler := NewUserHandler(context.Background(), sqlStore, cache)

	router := gin.New()
	router.GET("/user/:email", handler.GetUserByEmail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/john.doe@example.com", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["user"])
}

func TestUserHandler_GetUsers(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	mockCache, err := miniredis.Run()
	assert.NoError(t, err)
	cache := redis.NewClient(&redis.Options{
		Addr: mockCache.Addr(),
	})
	sqlStore := &store.UserStoreSql{Db: db}
	err = sqlStore.Migrate(context.Background())
	assert.NoError(t, err)

	handler := NewUserHandler(context.Background(), sqlStore, cache)

	users := []*store.UserSql{
		{
			FirstName: "John", 
			LastName: "Doe", 
			EmailAddress: "john.doe@example.com", 
			ParentUserId: 0,
		},
		{
			FirstName: "Jane", 
			LastName: "Smith", 
			EmailAddress: "jane.smith@example.com", 
			ParentUserId: 0,
		},
	}

	for _, user := range users{
		err = sqlStore.CreateUser(context.Background(), user) 
		assert.NoError(t, err) 
	}

	router := gin.New()
	router.GET("/users/:limit/:skip", handler.GetUsers) 

	w := httptest.NewRecorder() 
	req, _ := http.NewRequest("GET", "/users/10/0", nil) 
	router.ServeHTTP(w, req) 

	assert.Equal(t, http.StatusOK, w.Code) 

	var response map[string]interface{} 
	err = json.Unmarshal(w.Body.Bytes(), &response) 
	assert.NoError(t, err)  
	assert.NotNil(t, response["users"])  

}
