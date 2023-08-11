package wbs

import (
	"boilerplate-api/models"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn    *websocket.Conn
	Message chan *Message
	User    *models.User `json:"user"`
	Room    *Room        `json:"room"`
}

type Message struct {
	Content string       `json:"content"`
	RoomID  string       `json:"roomId"`
	User    *models.User `json:"user"`
}

func (c *Client) writeMessage() {
	fmt.Println("write message-->")
	defer func() {
		c.Conn.Close()
	}()
	for {
		message, ok := <-c.Message
		if !ok {
			fmt.Println("someting not okay--->")
			return
		}
		fmt.Println(&message, "write message,--->")
		c.Conn.WriteJSON(message)
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		fmt.Println("Reading message$$")
		_, m, err := c.Conn.ReadMessage()
		fmt.Println("this is read message")
		if err != nil {
			fmt.Println("Something went wrong")
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		fmt.Println(string(m), "MESSAGE FRON CLIENT")
		roomId := fmt.Sprintf("%v", c.Room.ID)
		fmt.Println(c.Room.ID, c.User.ID, "--------->roomId user Id")
		msg := &Message{
			Content: string(m),
			RoomID:  roomId,
			User:    c.User,
		}

		hub.Broadcast <- msg
		// h.Rooms[m.RoomID].Clients
	}
}
