package store 

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

)

type UserStoreSql struct{
	Db *gorm.DB 

}

type UserSql struct{
	//gorm.Model
	ID uint `gorm:"primaryKey" json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	EmailAddress string `json:"email_address"`
	MergedAt time.Time `json:"merged_at"`
	CreatedAt time.Time `json:"created_at"` 
	DeletedAt time.Time `json:"deleted_at"`
	ParentUserId uint `json:"parent_user_id"`
	ParentUser *UserSql `gorm:"foreignKey:ParentUserId"`
}

func (u *UserSql) TableName() string{
	return "users"
}

func (us *UserStoreSql) Migrate(ctx  context.Context) error {
	err := us.Db.AutoMigrate(&UserSql{}) 
	if err != nil{
		return errors.New("could not migrate the users table")
	}
	return nil 
}

func (us *UserStoreSql) CreateUser(ctx context.Context, user *UserSql) error{
	
	err := us.Db.Create(user).Error 
	if err != nil{
		return errors.New(fmt.Sprintf("could not create user an error occurred %v", err.Error()))
	}
	return nil 


}

func (us *UserStoreSql) GetUserByEmail(ctx context.Context, email string) (*UserSql, error){
	var user UserSql
	err  := us.Db.Where("email_address = ?", email).First(&user).Error 
	if err != nil{
		return nil, errors.New(fmt.Sprintf("user not found: %v", err))
	}

	return &user, nil 
}

func (us *UserStoreSql) GetUserById(ctx context.Context, userId uint) (*UserSql, error){
	var user UserSql 
	err := us.Db.Where("id = ?", userId).First(&user).Error 
	if err != nil{
		return nil, errors.New(fmt.Sprintf("user not found : %v", err)) 
	}

	return &user, nil 
}

func (us *UserStoreSql) GetUsers(ctx context.Context, limit int, skip int) ([]*UserSql, error){
	var users []*UserSql 
	err := us.Db.Limit(limit).Offset(skip).Find(&users).Error 
	if err != nil{
		return nil, errors.New(fmt.Sprintf("could not retrieve user"))
	}

	return users, nil 
}



