package models

//Login Model for users
type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	IsBlocked bool
}
