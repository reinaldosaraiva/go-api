package dto

type CreateProductDTO struct {
	Name string `json:"name"`
	Price int `json:"price"`
}

type CreateUserDTO struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginDTO struct {
	Email string `json:"email"`
	Password string `json:"password"`
}