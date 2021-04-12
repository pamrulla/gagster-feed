package models

import "gorm.io/gorm"

type Heart struct {
	gorm.Model
	Gag_Id  int `gorm:"<-:create;not null"`
	User_Id int `gorm:"<-:create;not null"`
}

type Hearts []Heart

//create a Heart
func CreateHeart(db *gorm.DB, hr *Heart) (err error) {
	err = db.Create(hr).Error
	if err != nil {
		return err
	}
	return nil
}

//get Hearts
func GetHearts(db *gorm.DB, hr *Hearts, gig_id string) (err error) {
	err = db.Where("gig_id = ?", gig_id).Find(hr).Error
	if err != nil {
		return err
	}
	return nil
}

//get Heart by id
func GetHeart(db *gorm.DB, hr *Heart, id string) (err error) {
	err = db.Where("id = ?", id).First(hr).Error
	if err != nil {
		return err
	}
	return nil
}

//update Heart
func UpdateHeart(db *gorm.DB, hr *Heart) (err error) {
	db.Save(hr)
	return nil
}

//delete Heart
func DeleteHeart(db *gorm.DB, hr *Heart, gag_id string, user_id string) (err error) {
	db.Where("gag_id = ? && user_id = ?", gag_id, user_id).Delete(hr)
	return nil
}

//empty Heart table
func EmptyHeartTable(db *gorm.DB) {
	db.Unscoped().Where("1 = 1").Delete(&Heart{})
}
