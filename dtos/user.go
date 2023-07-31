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
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Phone     string ` json:"phone" `
	Gender    string ` json:"gender"`
	Username  string ` json:"username"`
	IsPrivate bool   ` json:"is_private"`
}

type FollowUser struct {
	FollowUserId int64 `json:"follow_user_id" binding:"required"`
}

type FollowSuggestions struct {
	Id         int64  `json:"user_id"`
	Username   string `json:"username"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	IsFollowed bool   `json:"is_followed_by_them"`
}
