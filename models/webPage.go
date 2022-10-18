package models



type WebPage struct {
	Site       string `json:"site"`
	Num_Links  int    `json:"num_links"`
	Images     int    `json:"images"`
	Last_Fetch string `json:"last_fetch"`
}
