package store 

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserStoreSql_Migrate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	userStore := &UserStoreSql{Db: db}
	err = userStore.Migrate(context.Background())
	assert.NoError(t, err)

	// Check if the "users" table exists
	assert.True(t, db.Migrator().HasTable(&UserSql{}))
}

func TestUserStoreSql_CreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	userStore := &UserStoreSql{Db: db}
	userStore.Migrate(context.Background())

	user := &UserSql{
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
		ParentUserId: 0,
	}

	err = userStore.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	// Check if the user has been created
	var retrievedUser UserSql
	err = db.First(&retrievedUser).Error

	// Check if the record is not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		t.Error("expected user record to be found, but it wasn't")
	}
	assert.NoError(t, err)
	assert.Equal(t, user.FirstName, retrievedUser.FirstName)
	assert.Equal(t, user.LastName, retrievedUser.LastName)
	assert.Equal(t, user.EmailAddress, retrievedUser.EmailAddress)
}

func TestUserStoreSql_GetUserByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	userStore := &UserStoreSql{Db: db}
	userStore.Migrate(context.Background())

	// Insert a test user
	user := &UserSql{
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
		ParentUserId: 0,
	}

	err = db.Create(&user).Error
	assert.NoError(t, err)

	// Test get user by email
	retrievedUser, err := userStore.GetUserByEmail(context.Background(), "john.doe@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, user.EmailAddress, retrievedUser.EmailAddress)
}

func TestUserStoreSql_GetUserById(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	userStore := &UserStoreSql{Db: db}
	userStore.Migrate(context.Background())

	// Insert a user
	user := &UserSql{
		FirstName:    "John",
		LastName:     "Doe",
		EmailAddress: "john.doe@example.com",
		ParentUserId: 0,
	}

	err = db.Create(&user).Error
	assert.NoError(t, err)

	// Test GetUserById
	retrievedUser, err := userStore.GetUserById(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, user.ID, retrievedUser.ID)
}

func TestUserStore_GetUsers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	userStore := &UserStoreSql{Db: db}
	userStore.Migrate(context.Background())

	// Insert test users
	users := []*UserSql{
		{
			FirstName:    "John",
			LastName:     "Doe",
			EmailAddress: "john.doe@example.com",
			ParentUserId: 0,
		},
		{
			FirstName:    "Jane",
			LastName:     "Smith",
			EmailAddress: "jane.smith@example.com",
			ParentUserId: 0,
		},
	}

	for _, u := range users {
		err := db.Create(u).Error
		assert.NoError(t, err)
	}

	// Test GetUsers
	retrievedUsers, err := userStore.GetUsers(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Len(t, retrievedUsers, 2)
}

func TestUserSql_TableName(t *testing.T) {
	user := &UserSql{}
	assert.Equal(t, "users", user.TableName())
}
