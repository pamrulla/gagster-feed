package models

import "time"

type Gag struct {
	Id           int       `json:"Id"`
	User_Id      int       `json:"User_Id"`
	Path         string    `json:"Path"`
	Created_Date time.Time `json:"Created_Date"`
	Price        float32   `json:"Price"`
	IsEnabled    bool      `json:"IsEnabled"`
}

type Gags []Gag
