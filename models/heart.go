package models

import "time"

type Heart struct {
	Gag_Id   int       `json:"Gag_Id"`
	User_Id  int       `json:"User_Id"`
	Liked_On time.Time `json:"Liked_On"`
}

type Hearts []Heart
