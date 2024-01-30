package dto

type CreateProductRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
