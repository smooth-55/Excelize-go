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

func (c UserRepository) GetOneUser(Id int64) (models.User, map[string]interface{}, error) {

	var followers, following []*dtos.GetUserResponse
	var userModel models.User
	resp := make(map[string]interface{})

	err := c.db.DB.
		Model(&models.User{}).
		Where("id = ?", Id).
		First(&userModel).
		Error

	c.db.DB.Model(&models.FollowUser{}).
		Select("users.*").
		Joins("JOIN users on follow_user.followed_to_id = users.Id").
		Where("followed_by_id = ? AND is_approved = 1", Id).
		Find(&following)

	c.db.DB.Model(&models.FollowUser{}).
		Select("users.*").
		Joins("JOIN users on follow_user.followed_by_id = users.Id").
		Where("followed_to_id = ? AND is_approved = 1", Id).
		Find(&followers)

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

func (c UserRepository) FollowSuggestions(userId int64) ([]dtos.FollowSuggestions, error) {
	obj := []dtos.FollowSuggestions{}
	query := c.db.DB.Model(&models.User{}).
		Select(`
		users.id AS Id,
		users.username AS Username,
		users.full_name AS FullName,
		users.email AS Email,
		CASE WHEN follow_user.followed_to_id = ? THEN 1 ELSE 0 END AS IsFollowed
		`, userId).
		Joins("LEFT JOIN follow_user on users.id = follow_user.followed_by_id").
		Where("users.id NOT IN (select followed_to_id from follow_user where followed_by_id = ?) AND users.id != ?", userId, userId).
		Order("users.created_at DESC").
		Limit(5).
		Find(&obj)
	return obj, query.Error
}

func (c UserRepository) GetTwoWayFollowers(userId int64) ([]dtos.FollowSuggestions, error) {
	obj := []dtos.FollowSuggestions{}
	query := c.db.DB.Model(&models.FollowUser{}).
		Select(`
		follow_user.followed_to_id as Id, 
		usr.username as Username,
		usr.full_name as FullName,
		usr.email as Email,
		CASE WHEN usr.Id IS NOT NULL THEN 1 ELSE 0 END AS IsFollowed
		`).
		Joins("JOIN follow_user t2 ON follow_user.followed_to_id = t2.followed_by_id").
		Joins("JOIN users as usr on follow_user.followed_to_id = usr.id").
		Where("follow_user.followed_by_id <> t2.followed_by_id AND follow_user.followed_by_id = ?", userId).
		Group("usr.username").
		Find(&obj)
	return obj, query.Error
}
