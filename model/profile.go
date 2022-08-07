package model

type Profile struct {
	Country     string   `json:"country"`
	DisplayName string   `json:"display_name"`
	ID          string   `json:"id"`
	Images      []Images `json:"images"`
}

type Images struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}
