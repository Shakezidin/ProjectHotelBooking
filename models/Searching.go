package models

type SearchRequest struct {
	LocOrPlace       string `json:"locationorplace" binding:"required"`
	FromDate         string `json:"from_date" binding:"required" default:"current_date"`
	ToDate           string `json:"to_date" binding:"required" default:"tomorrow_date"`
	NumberOfChildren int    `json:"number_of_children" binding:"required" default:"0"`
	NumberOfAdults   int    `json:"number_of_adults" binding:"required" default:"1"`
}