package models

type Password struct {
	NewPassword     string `json:"new_password"`
	CurrentPassword string `json:"current_password"`
}
