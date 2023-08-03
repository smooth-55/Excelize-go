package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"gorm.io/gorm"
)

// UserRepository database structure
type RoomRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

// NewUserRepository creates a new User repository
func NewRoomRepository(db infrastructure.Database, logger infrastructure.Logger) RoomRepository {
	return RoomRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (c RoomRepository) WithTrx(trxHandle *gorm.DB) RoomRepository {
	if trxHandle == nil {
		c.logger.Zap.Error("Transaction Database not found in gin context. ")
		return c
	}
	c.db.DB = trxHandle
	return c
}

func (c RoomRepository) GetRoomById(Id int64) (models.Rooms, error) {
	var room models.Rooms
	return room, c.db.DB.Model(&models.Rooms{}).Where("id = ?", Id).
		First(&room).
		Error
}
