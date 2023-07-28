package dtos

// JWTLoginRequestData Request body data to authenticate user with jwt-auth
type JWTLoginRequestData struct {
	EmailOrUsername string `json:"email_or_username" validate:"required"`
	Password        string `json:"password" validate:"required"`
}
