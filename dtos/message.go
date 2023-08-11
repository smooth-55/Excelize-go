package dtos

import "boilerplate-api/models"

type AllConversations struct {
	ID       int                `json:"id"`
	UserId   int64              `json:"user_id"`
	Messages []*models.Messages `json:"messages"`
}
