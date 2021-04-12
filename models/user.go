package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	First_Name string `gorm:"<-:update;not null"`
	Last_Name  string `gorm:"<-:update;not null"`
	Email      string `gorm:"<-:create;not null"`
	Password   string `gorm:"<-:update;not null"`
	IsEnabled  bool   `gorm:"<-:update;default:0"`
}

type Users []User

//create a user
func CreateUser(db *gorm.DB, User *User) (err error) {
	err = db.Create(User).Error
	if err != nil {
		return err
	}
	return nil
}

//get users
func GetUsers(db *gorm.DB, User *Users) (err error) {
	err = db.Find(User).Error
	if err != nil {
		return err
	}
	return nil
}

//get user by id
func GetUser(db *gorm.DB, User *User, id string) (err error) {
	err = db.Where("ID = ?", id).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

//update user
func UpdateUser(db *gorm.DB, User *User) (err error) {
	db.Save(User)
	return nil
}

//delete user
func DeleteUser(db *gorm.DB, User *User, id string) (err error) {
	db.Where("ID = ?", id).Delete(User)
	return nil
}
