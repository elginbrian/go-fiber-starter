package domain

type UserRegistrationRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRegistrationResponse struct {
	Message string `json:"message"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}