package models

type Todo struct {
	Base
	Task         string `gorm:"column:task" json:"task" binding:"required"`
	Is_completed bool   `gorm:"column:is_completed" json:"is_completed"`
}

func (c Todo) TableName() string {
	return "Todo"
}
