package ws

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"api_kino/service/jwt_auth"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512

	ChannelUser              = ""
	ChannelSchedule          = "channel_schedule_"
	TopicPing                = "ping"
	TopicCreatedChat         = "created_chat"
	TopicCreatedNotification = "created_notification"
	TopicForceLogout         = "force_logout"
)

type Option struct {
	Code int64 `json:"code"`
}

type Request struct {
	Topic  string                 `json:"topic"`
	Option Option                 `json:"option"`
	Data   map[string]interface{} `json:"data"`
}

type Response struct {
	Topic  string      `json:"topic"`
	Option Option      `json:"option"`
	Data   interface{} `json:"data"`
}

type Channel struct {
	ID    string `json:"id"`
	Group string `json:"group"`
	Gin   *gin.Context
}

// connection is an middleman between the websocket connection and the H.
type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(gin *gin.Context) {
	ws, err := upgrader.Upgrade(gin.Writer, gin.Request, nil)
	if err != nil {
		log.Print(err)
		_ = ws.Close()
		return
	}
	jwt, err := jwt_auth.ValidateToken("Bearer " + gin.Param("auth"))
	if err != nil {
		_ = ws.Close()
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{gin, c, jwt.UserID}
	Hub.register <- s
	go s.writePump()
	go s.readPump()
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (s subscription) readPump() {
	c := s.conn
	defer func() {
		Hub.unregister <- s
		_ = c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				//log.Printf("error: %v", err)
			}
			break
		}
		var request Request
		if err := json.Unmarshal(msg, &request); err != nil {
			log.Printf("error: %v", err)
			continue
		}
		HandleData(&s, &request)
	}
}

func HandleData(s *subscription, request *Request) {
	switch request.Topic {
	case TopicPing:
		response := Response{
			Topic:  request.Topic,
			Option: Option{Code: 200},
		}
		BroadCastMessage(&response, &Channel{
			ID:    s.room,
			Group: ChannelUser,
			Gin:   s.gin,
		})
	case TopicCreatedChat:
		//db := database.DB
		//val, err := provider.CreateChat(db, s.room, request.Data)
		//if err != nil {
		//	return
		//}
		//response := Response{
		//	Topic:  request.Topic,
		//	Option: Option{Code: 200},
		//	Data:   val,
		//}
		//BroadCastMessage(&response, &Channel{
		//	ID:    val.FromID,
		//	Group: ChannelUser,
		//	Gin:   s.gin,
		//})
		//BroadCastMessage(&response, &Channel{
		//	ID:    val.ToID,
		//	Group: ChannelUser,
		//	Gin:   s.gin,
		//})
	default:
		response := Response{
			Topic:  request.Topic,
			Option: Option{Code: 200},
		}
		BroadCastMessage(&response, &Channel{
			ID:    s.room,
			Group: ChannelUser,
			Gin:   s.gin,
		})
	}
}

func BroadCastMessage(response *Response, channel *Channel) {
	data, err := json.Marshal(response)
	if err != nil {
		return
	}
	if channel.Group == ChannelSchedule {
		//db := database.DB
		//users, err := provider.GetUsersByCompany(db, channel.ID)
		//if err == nil {
		//	for _, v := range users {
		//		m := message{data, v.ID}
		//		Hub.broadcast <- m
		//	}
		//}
	} else {
		m := message{data, channel.ID}
		Hub.broadcast <- m
	}
}
