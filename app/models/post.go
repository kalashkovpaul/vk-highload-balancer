package models

//easyjson:json
type Posts []Post

//easyjson:json
type Post struct {
	Id       int64  `json:"id"`
	Parent   int64  `json:"parent"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited"`
	Forum    string `json:"forum"`
	Thread   int64  `json:"thread"`
	Created  string `json:"created"`
}

//easyjson:json
type PostUpdate struct {
	Message string `json:"message"`
}

//easyjson:json
type PostFull struct {
	Post   *Post   `json:"post"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}
