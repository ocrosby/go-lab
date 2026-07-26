package models

type Post struct {
	Id     int    `json:"Id"`
	Title  string `json:"Title"`
	Body   string `json:"Body"`
	UserId int    `json:"UserId"`
}

func NewPost() *Post {
	return new(Post)
}
