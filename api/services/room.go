package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"
	"gorm.io/gorm"
)

type RoomService struct {
	repository repository.RoomRepository
}

// NewUserService Creates New user service
func NewRoomService(repository repository.RoomRepository) RoomService {
	return RoomService{
		repository: repository,
	}
}

// WithTrx repository with transaction
func (c RoomService) WithTrx(trxHandle *gorm.DB) RoomService {
	c.repository = c.repository.WithTrx(trxHandle)
	return c
}

func (c RoomService) GetRoomById(Id int64) (models.Rooms, error) {
	return c.repository.GetRoomById(Id)
}
