package models

import "time"

//easyjson:json
type Threads []Thread

//easyjson:json
type Thread struct {
	Id      int64     `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int32     `json:"votes"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
}

//easyjson:json
type ThreadUpdate struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}
