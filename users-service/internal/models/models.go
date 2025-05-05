package models

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type NewUserRequest struct {
	Email string `json:"email"`
}
