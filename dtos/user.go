package dtos

import "boilerplate-api/models"

// CreateUserRequestData Request body data to create user
type CreateUserRequestData struct {
	Email    string `gorm:"email" json:"email" validate:"required,email"`
	FullName string `gorm:"full_name" json:"full_name" validate:"required"`
	Password string `gorm:"password" json:"password" validate:"required"`
	Username string `gorm:"username" json:"username" validate:"required"`
}

func (c *CreateUserRequestData) GetUserObj() models.User {
	return models.User{
		Email:    c.Email,
		FullName: c.FullName,
		Password: c.Password,
		Username: c.Username,
	}
}

// GetUserResponse Dtos for User model
type GetUserResponse struct {
	models.User
	Password string
}
