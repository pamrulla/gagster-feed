package models

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	First_Name string `gorm:"not null"`
	Last_Name  string `gorm:"not null"`
	Email      string `gorm:"<-:create;not null"`
	Password   string `gorm:"not null"`
	IsEnabled  bool   `gorm:"<-:update;default:0"`
}

type Users []User

type UserHandlerInterface interface {
	CreateUser(db *gorm.DB, User *User) (err error)
	GetUsers(db *gorm.DB, User *Users) (err error)
	GetUser(db *gorm.DB, User *User, id string) (err error)
	UpdateUser(db *gorm.DB, User *User) (err error)
	DeleteUser(db *gorm.DB, User *User, id string) (err error)
	EmptyUserTable(db *gorm.DB)
	LoginUser(db *gorm.DB, User *User, email string, password string) (err error)
	IsUserExists(db *gorm.DB, email string) bool
}

type UserHandler struct{}

//create a user
func (u UserHandler) CreateUser(db *gorm.DB, User *User) (err error) {
	err = db.Create(User).Error
	if err != nil {
		return err
	}
	return nil
}

//get users
func (u UserHandler) GetUsers(db *gorm.DB, User *Users) (err error) {
	err = db.Find(User).Error
	if err != nil {
		return err
	}
	return nil
}

//get user by id
func (u UserHandler) GetUser(db *gorm.DB, User *User, id string) (err error) {
	err = db.Where("ID = ?", id).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

//is user exists
func (u UserHandler) IsUserExists(db *gorm.DB, email string) bool {
	var usr User
	fmt.Println(email)
	err := db.Where("email = ?", email).First(&usr).Error
	if err == nil {
		return !false
	}
	return !true
}

//log in user by email and password
func (u UserHandler) LoginUser(db *gorm.DB, User *User, email string, password string) (err error) {
	err = db.Where("email = ? AND password = ?", email, password).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

//update user
func (u UserHandler) UpdateUser(db *gorm.DB, User *User) (err error) {
	db.Save(User)
	return nil
}

//delete user
func (u UserHandler) DeleteUser(db *gorm.DB, User *User, id string) (err error) {
	db.Where("ID = ?", id).Delete(User)
	return nil
}

//empty User table
func (u UserHandler) EmptyUserTable(db *gorm.DB) {
	db.Unscoped().Where("1 = 1").Delete(&User{})
}
