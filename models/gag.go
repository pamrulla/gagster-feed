package models

import (
	"gorm.io/gorm"
)

type Gag struct {
	gorm.Model
	User_Id   int     `gorm:"<-:create;not null"`
	Path      string  `gorm:"<-:update;not null"`
	Price     float32 `gorm:"<-:update;not null"`
	IsEnabled bool    `gorm:"<-:update;default:1"`
}

type Gags []Gag

//create a Gag
func CreateGag(db *gorm.DB, gag *Gag) (err error) {
	err = db.Create(gag).Error
	if err != nil {
		return err
	}
	return nil
}

//get Gags
func GetGags(db *gorm.DB, gag *Gags, user_id string) (err error) {
	err = db.Where("user_id = ?", user_id).Find(gag).Error
	if err != nil {
		return err
	}
	return nil
}

//get Gag by id
func GetGag(db *gorm.DB, gag *Gag, id string) (err error) {
	err = db.Where("id = ?", id).First(gag).Error
	if err != nil {
		return err
	}
	return nil
}

//update Gag
func UpdateGag(db *gorm.DB, gag *Gag) (err error) {
	db.Save(gag)
	return nil
}

//delete Gag
func DeleteGag(db *gorm.DB, gag *Gag, id string) (err error) {
	db.Where("id = ?", id).Delete(gag)
	return nil
}

//empty Gag table
func EmptyGagTable(db *gorm.DB) {
	db.Unscoped().Where("1 = 1").Delete(&Gag{})
}
