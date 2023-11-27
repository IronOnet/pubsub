package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/irononet/consumer/store"
)

type UserHandler struct {
	Ctx   context.Context
	Cache *redis.Client
	Store *store.UserStoreSql
}

func NewUserHandler(ctx context.Context, store *store.UserStoreSql, cache *redis.Client) *UserHandler {
	return &UserHandler{
		Ctx:   ctx,
		Cache: cache,
		Store: store,
	}
}

// TODO: Write Swagger Specification Above
func (u *UserHandler) GetOneUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Println("could not parse user id: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	idUint := uint(id)

	// Try to get user from cache
	cacheKey := fmt.Sprintf("user:%d", id)
	cachedUser, err := u.getUserFromCache(cacheKey)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"user": cachedUser})
		return
	}

	// If not found in cache, fetch from the store
	user, err := u.Store.GetUserById(context.Background(), idUint)
	if err != nil {
		log.Println("error getting user by id:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	// Save the user in cache for future requests
	if err := u.setUserInCache(cacheKey, user); err != nil {
		log.Println("error saving user in cache: ", err)
	}

	c.JSON(http.StatusOK, gin.H{"user": user})

}

// TODO: Write Swagger Specification Above
func (u *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	cacheKey := fmt.Sprintf("user:%s", email)
	cachedUser, err := u.getUserFromCache(cacheKey)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"user": cachedUser})
		return
	}

	user, err := u.Store.GetUserByEmail(context.Background(), email)
	if err != nil {
		log.Println("error getting user by email: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	// Save use in cache for future requests
	if err := u.setUserInCache(cacheKey, user); err != nil {
		log.Println("error saving user in cache:", err)
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// TODO: Write Swagger Specification Above
func (u *UserHandler) GetUsers(c *gin.Context) {
	limitStr := c.Param("limit")
	skipStr := c.Param("skip")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		log.Println("error converting limit to integer:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	skip, err := strconv.ParseInt(skipStr, 10, 64)
	if err != nil {
		log.Println("error converting skip to integer:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skip"})
		return
	}

	limitInt := int(limit)
	skipInt := int(skip)

	cacheKey := "all_users"
	cachedUsers, err := u.getUserFromCache(cacheKey)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"users": cachedUsers})
		return
	}

	// If not found in cache, fetch from the store
	users, err := u.Store.GetUsers(context.Background(), limitInt, skipInt)
	if err != nil {
		log.Println("error getting users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (u *UserHandler) getUserFromCache(key string) (*store.UserSql, error) {
	val, err := u.Cache.Get(key).Result()
	if err != nil {
		return nil, err
	}

	// Unmarshall cached value into User struct
	var user store.UserSql
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserHandler) setUserInCache(key string, user *store.UserSql) error {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshalling user: %v", err))
	}

	return u.Cache.Set(key, userJSON, 1*time.Hour).Err()
}
