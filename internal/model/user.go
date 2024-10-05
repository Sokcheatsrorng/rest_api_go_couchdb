package model


type User struct {

	ID string `json:"_id"`

	Rev string `json:"_rev"`

	Username string `json:"username"`

	Password string `json:"password"`

	Email string `json:"email"`

}