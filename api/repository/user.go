package repository

import (
	"boilerplate-api/dtos"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/paginations"

	"gorm.io/gorm"
)

// UserRepository database structure
type UserRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

// NewUserRepository creates a new User repository
func NewUserRepository(db infrastructure.Database, logger infrastructure.Logger) UserRepository {
	return UserRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (c UserRepository) WithTrx(trxHandle *gorm.DB) UserRepository {
	if trxHandle == nil {
		c.logger.Zap.Error("Transaction Database not found in gin context. ")
		return c
	}
	c.db.DB = trxHandle
	return c
}

// Create user
func (c UserRepository) Create(User models.User) error {
	return c.db.DB.Create(&User).Error
}

// GetAllUsers Get All users
func (c UserRepository) GetAllUsers(pagination paginations.UserPagination) (users []dtos.GetUserResponse, count int64, err error) {
	queryBuilder := c.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset).Order("created_at desc")
	queryBuilder = queryBuilder.Model(&models.User{})

	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		queryBuilder.Where(c.db.DB.Where("`users`.`name` LIKE ?", searchQuery))
	}

	return users, count, queryBuilder.
		Find(&users).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error
}

func (c UserRepository) GetOneUser(Id string) (dtos.GetUserResponse, map[string]interface{}, error) {

	followers := []*models.FollowUser{}
	following := []*models.FollowUser{}
	userModel := dtos.GetUserResponse{}
	resp := make(map[string]interface{})

	err := c.db.DB.
		Model(&userModel).
		Where("id = ?", Id).
		First(&userModel).
		Error

	c.db.DB.Model(&models.FollowUser{}).
		Where("user_id = ? AND is_approved = 1", Id).
		Find(&followers)

	c.db.DB.Model(&models.FollowUser{}).
		Where("followed_user_id = ? AND is_approved = 1", Id).
		Find(&following)

	resp["followers"] = followers
	resp["following"] = following

	if err != nil {
		return userModel, resp, err
	}
	return userModel, resp, err
}

func (c UserRepository) GetOneUserWithEmail(Email string) (user models.User, err error) {
	return user, c.db.DB.Model(&user).
		Where("email = ?", Email).
		First(&user).
		Error
}

func (c UserRepository) GetOneUserWithUsername(username string) (user models.User, err error) {
	return user, c.db.DB.Model(&user).
		Where("username = ?", username).
		First(&user).
		Error
}

func (c UserRepository) GetOneUserWithPhone(Phone string) (user models.User, err error) {
	return user, c.db.DB.
		First(&user, "phone = ?", Phone).
		Error

}

func (c UserRepository) FollowUser(obj models.FollowUser) error {
	return c.db.DB.Create(&obj).Error
}
