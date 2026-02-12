package requests

type UserRegisterRequest struct {
	Login           string `json:"login" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type UserLoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserJWTRefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required,min=10"`
}
