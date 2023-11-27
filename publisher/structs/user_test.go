package structs

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

func TestUser(t *testing.T){
	// Create a simple user instance 
	tm , err := strconv.ParseInt("136145982200", 10, 64)
	if err != nil{
		panic(err) 
	}
	user := User{
		ID: 1, 
		FirstName: "john", 
		LastName: "maynard", 
		EmailAddress: "johnmaynard@biz.com", 
		CreatedAt: time.Unix(tm, 0), 
		DeletedAt: time.Unix(tm, 0), 
		MergedAt: time.Unix(tm, 0), 
		ParentUserId: 0,
	}

	jsonData, err := json.Marshal(user) 
	if err != nil{
		t.Errorf("error marshalling user to JSON: %v", err) 
	}

	// Convert the JSON back to a user struct 
	var parsedUser User 
	err = json.Unmarshal(jsonData, &parsedUser) 
	if err != nil{
		t.Errorf("error unmarshalling JSON to user: %v", err)
	}

	// Compare the original user and the parsed user 
	if user != parsedUser{
		t.Errorf("parsed user does not match the original user")
	}
}