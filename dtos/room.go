package dtos

type CreateRoomDTO struct {
	Id   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
