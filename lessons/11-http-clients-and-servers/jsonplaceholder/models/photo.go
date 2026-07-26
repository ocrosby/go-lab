package models

type Photo struct {
	AlbumId  int    `json:"albumId"`
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	ThumbUrl string `json:"thumbnailUrl"`
}
