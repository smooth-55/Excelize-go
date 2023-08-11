package wbs

import (
	"boilerplate-api/models"
	"fmt"
)

type Room struct {
	*models.Rooms
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
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
