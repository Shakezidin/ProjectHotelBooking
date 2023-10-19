package models

//SearchRequest for binding seraching details
type SearchRequest struct {
	LocOrPlace       string `json:"locationorplace" binding:"required"`
	FromDate         string `json:"from_date" binding:"required" default:"current_date"`
	ToDate           string `json:"to_date" binding:"required" default:"tomorrow_date"`
	NumberOfChildren int    `json:"number_of_children" default:"0"`
	NumberOfAdults   int    `json:"number_of_adults" default:"1"`
}
