package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"

	"gorm.io/gorm"
)

type MessageService struct {
	repository repository.MessageRepository
}

// NewUserService Creates New user service
func NewMessageService(repository repository.MessageRepository) MessageService {
	return MessageService{
		repository: repository,
	}
}

// WithTrx repository with transaction
func (c MessageService) WithTrx(trxHandle *gorm.DB) MessageService {
	c.repository = c.repository.WithTrx(trxHandle)
	return c
}

func (c MessageService) GetMyConversations(userId int64) ([]models.RoomUsers, int64, error) {
	return c.repository.GetMyConversations(userId)
}

func (c MessageService) GetOneRoomById(rommId int64) (models.Rooms, error) {
	return c.repository.GetOneRoomById(rommId)
}
