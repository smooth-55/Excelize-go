package wbs

import (
	"boilerplate-api/api/services"
	"boilerplate-api/models"
	"fmt"
	"strconv"
	"time"
)

type Room struct {
	*models.Rooms
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms          map[string]*Room
	Register       chan *Client
	Unregister     chan *Client
	Broadcast      chan *Message
	messageService services.MessageService
}

func NewHub(
	messageService services.MessageService,
) *Hub {
	return &Hub{
		Rooms:          make(map[string]*Room),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Broadcast:      make(chan *Message, 5),
		messageService: messageService,
	}
}

func (h *Hub) RegisterClient(cl *Client) {
	userId := fmt.Sprintf("%v", cl.User.ID)
	roomId := fmt.Sprintf("%v", cl.Room.ID)
	if _, ok := h.Rooms[roomId]; ok {
		r := h.Rooms[roomId]
		if _, ok := r.Clients[userId]; !ok {
			r.Clients[userId] = cl
		}
	}
}

func (h *Hub) UnRegisterClient(cl *Client) {
	userId := fmt.Sprintf("%v", cl.User.ID)
	roomId := fmt.Sprintf("%v", cl.Room.ID)
	if _, ok := h.Rooms[roomId]; ok {
		if _, ok := h.Rooms[roomId].Clients[userId]; ok {
			delete(h.Rooms[roomId].Clients, userId)
			close(cl.Message)
		}
	}
}

func (h *Hub) BroadcastMessage(msg *Message) {
	if _, ok := h.Rooms[msg.RoomID]; ok {
		for _, cl := range h.Rooms[msg.RoomID].Clients {
			cl.Message <- msg
			h.SaveMessage(msg)
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			h.RegisterClient(cl)
		case cl := <-h.Unregister:
			h.UnRegisterClient(cl)
		case m := <-h.Broadcast:
			h.BroadcastMessage(m)
		}
	}
}

func (h *Hub) SaveMessage(msg *Message) {
	_intRoom, _ := strconv.ParseInt(msg.RoomID, 10, 64)
	var obj models.Messages
	obj.RoomId = _intRoom
	obj.Message = msg.Content
	obj.MessageBy = msg.User.ID
	obj.CreatedAt = time.Now()
	if err := h.messageService.CreateMessage(obj); err != nil {
		fmt.Println("Err while saving to db---> [CreateMessage] ", err)
	}
}
