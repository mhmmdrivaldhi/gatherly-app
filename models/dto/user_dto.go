package dto

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Age      int    `json:"age" validate:"required,gte=0,lte=150"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"omitempty,oneof=user admin"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" validate:"omitempty,min=3,max=100"`
	Age      int    `json:"age" validate:"omitempty,gte=0,lte=150"`
	Username string `json:"username" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email" validate:"omitempty,email,max=100"`
	Role     string `json:"role" validate:"omitempty,oneof=user admin"`
}

type LoginRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LoginResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	Latitude  float64      `json:"latitude,omitempty"`
	Longitude float64      `json:"longitude,omitempty"`
}
