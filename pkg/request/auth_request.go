package request

type UserRegistrationRequest struct {
	Username string `json:"username" validate:"required" example:"john_doe"`
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"P@ssw0rd123"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"P@ssw0rd123"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=6" example:"OldP@ssw0rd"`
	NewPassword string `json:"new_password" validate:"required,min=6" example:"NewP@ssw0rd123"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"your_refresh_token_here"`
}