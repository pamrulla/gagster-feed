package models

type User struct {
	Id           int    `json:"Id"`
	First_Name   string `json:"First_Name"`
	Last_Name    string `json:"Last_Name"`
	Email        string `json:"Email"`
	Password     string `json:"Password"`
	Created_Date string `json:"Created_Date"`
}

type Users []User
