package model

type Job struct {
	Title         string   `json:"title"`
	Location      string   `json:"location"`
	Link          string   `json:"link"`
	Qualification []string `json:"qualification"`
}
