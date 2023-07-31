package models

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Base
	Email     string `gorm:"column:email" json:"email" validate:"required,email"`
	FullName  string `gorm:"column:full_name" json:"full_name" validate:"required"`
	Phone     string `gorm:"column:phone;default:null" json:"phone" validate:"required,phone"`
	Gender    string `gorm:"column:gender" json:"gender"`
	Password  string `gorm:"column:password" json:"-"`
	Username  string `gorm:"column:username" json:"username" validate:"required"`
	IsPrivate bool   `gorm:"column:is_private;default:0" json:"is_private"`
}

// TableName gives table name of model
func (u *User) TableName() string {
	return "users"
}

// ToMap convert User to map
func (u *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"email":      u.Email,
		"full_name":  u.FullName,
		"phone":      u.Phone,
		"username":   u.Username,
		"gender":     u.Gender,
		"is_private": u.IsPrivate,
	}
}

// BeforeCreate Runs before inserting a row into table
func (u *User) BeforeCreate(db *gorm.DB) error {
	var Zap *zap.SugaredLogger
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.Password = string(password)
	if err != nil {
		Zap.Error("Error decrypting plain password to hash", err.Error())
		return err
	}
	return nil
}

type FollowUser struct {
	Base
	FollowedById int64 `gorm:"column:followed_by_id" json:"followed_by_id" binding:"-"`
	FollowedToId int64 `gorm:"column:followed_to_id" json:"followed_to_id"`
	IsApproved   bool  `gorm:"column:is_approved;default:1" json:"is_approved"`
}

func (u FollowUser) TableName() string {
	return "follow_user"
}
