package dtos

import "boilerplate-api/models"

type AllConversations struct {
	UserId   int64              `json:"user_id"`
	Messages []*models.Messages `json:"messages"`
}
