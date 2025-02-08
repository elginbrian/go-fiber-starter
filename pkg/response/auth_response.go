package response

type LoginResponse struct {
	Status string    `json:"status" example:"success"`
	Data   LoginData `json:"data"`
	Code   int       `json:"code" example:200`
}

type LoginData struct {
	AccessToken  string `json:"access_token" example:"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type RegisterResponse struct {
	Status string       `json:"status" example:"success"`
	Data   RegisterData `json:"data"`
	Code   int          `json:"code" example:201`
}

type RegisterData struct {
	Message string `json:"message" example:"User registered successfully"`
}

type GetCurrentUserResponse struct {
	Status string `json:"status" example:"success"`
	Data   User   `json:"data"`
	Code   int    `json:"code" example:200`
}

type ChangePasswordResponse struct {
	Status string             `json:"status" example:"success"`
	Data   ChangePasswordData `json:"data"`
	Code   int                `json:"code" example:200`
}

type ChangePasswordData struct {
	Message string `json:"message" example:"Password changed successfully"`
}

type RefreshTokenResponse struct {
	Status string           `json:"status" example:"success"`
	Data   RefreshTokenData `json:"data"`
	Code   int              `json:"code" example:200`
}

type RefreshTokenData struct {
	AccessToken string `json:"access_token" example:"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}