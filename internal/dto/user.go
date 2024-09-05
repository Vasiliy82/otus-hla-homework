package dto

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterUserRequest struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Birthdate  string `json:"birthdate"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}
