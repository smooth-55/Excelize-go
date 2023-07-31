package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TodoRepository database structure
type TodoRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

// NewTodoRepository creates a new Todo repository
func NewTodoRepository(db infrastructure.Database, logger infrastructure.Logger) TodoRepository {
	return TodoRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (c TodoRepository) WithTrx(trxHandle *gorm.DB) TodoRepository {
	if trxHandle == nil {
		c.logger.Zap.Error("Transaction Database not found in gin context. ")
		return c
	}
	c.db.DB = trxHandle
	return c
}

// Create Todo
func (c TodoRepository) Create(Todo models.Todo) (models.Todo, error) {
	return Todo, c.db.DB.Create(&Todo).Error
}

// Create Todo
func (c TodoRepository) BulkCreateTodo(Todo []*models.Todo) error {
	return c.db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		UpdateAll: true,
	}).Create(&Todo).Error
}

// Create Todo

// GetAllTodo Get All todos
func (c TodoRepository) GetAllTodo() ([]models.Todo, int64, error) {
	var todos []models.Todo
	var totalRows int64 = 0
	queryBuilder := c.db.DB.Model(&models.Todo{})
	err := queryBuilder.
		Find(&todos).
		Offset(-1).
		Limit(-1).
		Count(&totalRows).Error
	return todos, totalRows, err
}

// GetOneTodo Get One Todo By Id
func (c TodoRepository) GetOneTodo(ID int64) (models.Todo, error) {
	Todo := models.Todo{}
	return Todo, c.db.DB.
		Where("id = ?", ID).First(&Todo).Error
}
