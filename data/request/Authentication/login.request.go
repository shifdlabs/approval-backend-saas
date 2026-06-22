package authentication

type LogInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1,max=200"`
}

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
