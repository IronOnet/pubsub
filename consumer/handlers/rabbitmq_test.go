package handlers

import (
	"context"
	//"encoding/json"
	//"errors"
	"testing"
	//"time"

	"github.com/irononet/consumer/store"
	//"github.com/streadway/amqp"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserStore struct{
	mock.Mock
}

func (m *MockUserStore) CreateUser(ctx context.Context, user *store.UserSql) error{
	args := m.Called(ctx, user) 
	return args.Error(0) 
}

func TestInitChannel(t *testing.T){
	// mockUserStore := new(MockUserStore)  

	// // Mock the CreateUser method to return a predefined error 
	// expectedError := errors.New("method CreateUser error") 
	// mockUserStore.On("CreateUser", mock.Anything, mock.Anything).Return(expectedError) 

	// // Create a mock rabbitmq channel 
	// mockChannel := &amqp.Channel{} 
	// mockChannel.
}