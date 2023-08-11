package wbs

import (
	"boilerplate-api/api/services"
	"boilerplate-api/infrastructure"
	"boilerplate-api/responses"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub            *Hub
	userService    services.UserService
	logger         infrastructure.Logger
	messageService services.MessageService
	jwtService     services.JWTAuthService
	env            infrastructure.Env
}

func NewHandler(
	h *Hub,
	userService services.UserService,
	logger infrastructure.Logger,
	messageService services.MessageService,
	jwtService services.JWTAuthService,
	env infrastructure.Env,
) *Handler {
	return &Handler{
		hub:            h,
		userService:    userService,
		logger:         logger,
		messageService: messageService,
		jwtService:     jwtService,
		env:            env,
	}
}

// type CreateRoomReq struct {
// 	ID   string `json:"id"`
// 	Name string `json:"name"`
// }

// func (h *Handler) CreateRoom(c *gin.Context) {
// 	var req CreateRoomReq
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	h.hub.Rooms[req.ID] = &Room{
// 		ID:      req.ID,
// 		Name:    req.Name,
// 		Clients: make(map[string]*Client),
// 	}

// 	c.JSON(http.StatusOK, req)
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoomTest(c *gin.Context) {
	responses.JSON(c, http.StatusOK, "Without authentication")
}

type AuthStruct struct {
	Authorization string `json:"Authorization"`
}

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomID := c.Param("roomId")
	clientID := c.Query("user_id")
	intUserId, _ := strconv.ParseInt(clientID, 10, 64)
	intRoomId, _ := strconv.ParseInt(roomID, 10, 64)
	getOneUser, _, _ := h.userService.GetOneUser(intUserId)
	getOneRoom, _ := h.messageService.GetOneRoomById(intRoomId)

	if _, ok := h.hub.Rooms[roomID]; !ok {
		h.hub.Rooms[roomID] = &Room{
			Rooms:   &getOneRoom,
			Clients: make(map[string]*Client),
		}
	}

	if getOneUser.ID != 0 {
		h.logger.Zap.Info("Client exists")
		cl := &Client{
			Conn:    conn,
			Message: make(chan *Message, 10),
			Room: &Room{
				Rooms: &getOneRoom,
			},
			User: &getOneUser,
		}
		// messageContent := fmt.Sprintf("User %v has joined the room", getOneUser.Username)
		// m := &Message{
		// 	Content: messageContent,
		// 	RoomID:  roomID,
		// 	User:    &getOneUser,
		// }
		h.hub.Register <- cl
		// h.hub.Broadcast <- m

		go cl.writeMessage()
		cl.readMessage(h.hub)
	}

}
