package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"
)

// TodoService struct
type TodoService struct {
	repository repository.TodoRepository
}

// NewTodoService creates a new TodoService
func NewTodoService(repository repository.TodoRepository) TodoService {
	return TodoService{
		repository: repository,
	}
}

// CreateTodo call to create the Todo
func (c TodoService) CreateTodo(todo models.Todo) (models.Todo, error) {
	return c.repository.Create(todo)
}

// GetAllTodo call to create the Todo
func (c TodoService) GetAllTodo() ([]models.Todo, int64, error) {
	return c.repository.GetAllTodo()
}

// GetOneTodo Get One Todo By Id
func (c TodoService) GetOneTodo(ID int64) (models.Todo, error) {
	return c.repository.GetOneTodo(ID)
}
