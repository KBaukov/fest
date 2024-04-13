package ent

import "time"

type UserAuthDataset struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

type UserRegDataset struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Taun_Number string `json:"taun_number"`
	FormId      string `json:"formid"`
	FormName    string `json:"formname"`
	TranId      string `json:"tranid"`
	COOKIES     string `json:"COOKIES"`
}

type User struct {
	ID         int32     `json:"id"`
	Email      string    `json:"email"`
	Pass       string    `json:"pass"`
	ActiveFlag bool      `json:"active_flag"`
	LastVisit  time.Time `json:"last_visit"`
}
