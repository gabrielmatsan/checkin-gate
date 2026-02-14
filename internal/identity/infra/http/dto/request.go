package dto

type GoogleCallbackRequest struct {
	Code  string `validate:"required"`
	State string `validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
