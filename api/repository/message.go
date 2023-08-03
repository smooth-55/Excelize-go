package repository

import (
	"boilerplate-api/dtos"
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

func (c MessageRepository) GetMyConversations(userId int64) ([]dtos.AllConversations, int64, error) {
	var resp []dtos.AllConversations
	var count int64
	query := c.db.DB.Model(&models.Rooms{}).
		Where("message_from_id = ? OR message_to_id = ?", userId, userId).
		Find(&resp).
		Count(&count)
	return resp, count, query.Error
}
