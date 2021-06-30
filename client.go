package main
import(
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)
type Data struct {
	Ip       string   `json:"ip"`
	//用户名
	User     string   `json:"user"`
	//消息发送者
	From     string   `json:"from"`
	Type     string   `json:"type"`
	Content  string   `json:"content"`
	UserList []string `json:"user_list"`
	RoomId string `json:"room_id"'`
}
type connection struct {
	ws           *websocket.Conn
	message_chan chan []byte
	data         *Data
	//roomID     []byte
}
var WsUpgrader = &websocket.Upgrader{
	ReadBufferSize: 512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool { return true },
}
func myws(context *gin.Context){
	ws,err:= WsUpgrader.Upgrade(context.Writer,context.Request,nil)
	if err!=nil{
		return
	}
	c:=&connection{message_chan: make(chan []byte, 256), ws: ws, data: &Data{}}
	h.register <- c
	go c.writer()
	c.reader()
	defer func() {
		c.data.Type = "logout"
		user_list = deleteUser(user_list, c.data.User)
		c.data.UserList = user_list
		c.data.Content = c.data.User
		data_b, _ := json.Marshal(c.data)
		h.broadcast <- data_b
		h.unregister <- c
	}()
}
func (c *connection) writer() {
	for message := range c.message_chan {
		c.ws.WriteMessage(websocket.TextMessage, message)
	}
	c.ws.Close()
}

var user_list = []string{}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			h.unregister <- c
			break
		}
		json.Unmarshal(message, &c.data)
		switch c.data.Type {
		case "login":
			c.data.User = c.data.Content
			c.data.From = c.data.User
			user_list = append(user_list, c.data.User)
			c.data.UserList = user_list
			fmt.Println("data is",c.data.RoomId,c.data.User)
			h.roomID[c]=c.data.RoomId
			data_b, _ := json.Marshal(c.data)
			//h.register <- c
			h.broadcast<-data_b
		case "user":
			c.data.Type = "user"
			data_b, _ := json.Marshal(c.data)
			h.broadcast <- data_b
		case "logout":
			c.data.Type = "logout"
			user_list = deleteUser(user_list, c.data.User)
			data_b, _ := json.Marshal(c.data)
			h.broadcast <- data_b
			h.unregister <- c
		default:
			fmt.Print("========default================")
		}
	}
}

func deleteUser(slice []string, user string) []string {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var n_slice = []string{}
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			n_slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	fmt.Println(n_slice)
	return n_slice
}
