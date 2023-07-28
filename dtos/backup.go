package dtos

import "boilerplate-api/models"

type BackupDTO struct {
	Users []*models.User `json:"users"`
	Todo  []*models.Todo `json:"todo"`
}

type FilePath struct {
	Path     string `json:"path" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}
