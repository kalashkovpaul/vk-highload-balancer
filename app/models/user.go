package models

//easyjson:json
type Users []User

//easyjson:json
type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

//easyjson:json
type UserUpdate struct {
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}
