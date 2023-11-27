package structs

import "time"



type User struct{
	ID int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
	MergedAt time.Time `json:"merged_at"`
	ParentUserId int `json:"parent_user_id"`
}