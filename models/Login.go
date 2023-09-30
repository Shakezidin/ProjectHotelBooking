package models

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Is_block bool
}
