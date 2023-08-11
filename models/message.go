package models

type PhoneMessage struct {
	Phone   string
	Message string
}

type Rooms struct {
	Base
	Name  string       `gorm:"column:name" json:"name" binding:"required"`
	Users []*RoomUsers `gorm:"foreignKey:room_id" json:"users"`
}

func (c Rooms) TableName() string {
	return "rooms"
}

type Messages struct {
	Base
	RoomId    int64  `gorm:"column:message_room_id" json:"message_room_id" binding:"required"`
	MessageBy int64  `gorm:"column:message_by" json:"message_by" binding:"required"`
	Message   string `gorm:"column:message" json:"message" binding:"required"`
}

func (c Messages) TableName() string {
	return "messages"
}

type RoomUsers struct {
	Base
	RoomId int64 `gorm:"column:room_id" json:"room_id" binding:"required"`
	Room   Rooms `gorm:"foreignKey:room_id" json:"room"`
	UserId int64 `gorm:"column:user_id" json:"user_id" binding:"required"`
	User   User  `gorm:"foreignKey:user_id" json:"user"`
}

func (c RoomUsers) TableName() string {
	return "room_users"
}
