package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"

	"gorm.io/gorm"
)

// UserRepository database structure
type MessageRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

// NewUserRepository creates a new User repository
func NewMessageRepository(db infrastructure.Database, logger infrastructure.Logger) MessageRepository {
	return MessageRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (c MessageRepository) WithTrx(trxHandle *gorm.DB) MessageRepository {
	if trxHandle == nil {
		c.logger.Zap.Error("Transaction Database not found in gin context. ")
		return c
	}
	c.db.DB = trxHandle
	return c
}

func (c MessageRepository) GetMyConversations(userId int64) ([]models.RoomUsers, int64, error) {
	var resp []models.RoomUsers
	var count int64
	query := c.db.DB.Model(&models.RoomUsers{}).
		Preload("Room").
		Preload("User").
		Where("user_id = ? ", userId).
		Find(&resp).
		Count(&count)
	return resp, count, query.Error
}

func (c MessageRepository) GetOneRoomById(roomId int64) (models.Rooms, error) {
	var resp models.Rooms
	query := c.db.DB.Model(&models.Rooms{}).
		Preload("Users").
		Preload("Users.User").
		Where("id = ? ", roomId).
		Find(&resp)
	return resp, query.Error
}

func (c MessageRepository) CreateMessage(msg models.Messages) error {
	return c.db.DB.Create(&msg).Error
}

func (c MessageRepository) GetAllMessagesByRoomId(roomId int64) ([]*models.Messages, error) {
	var msgs []*models.Messages
	return msgs, c.db.DB.Model(&models.Messages{}).
		Preload("User").
		Where("room_id = ?", roomId).
		Order("created_at").
		Find(&msgs).
		Error
}
