package main

import (
	"encoding/json"
	"fmt"
)

var h = hub{
	connections: make(map[*connection]bool),
	unregister:  make(chan *connection),
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	roomID: make(map[*connection]string),
}

type hub struct {
	connections map[*connection]bool
	broadcast   chan []byte
	register    chan *connection
	unregister  chan *connection
	roomID     map[*connection]string
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			//_,tmp,_:=c.ws.ReadMessage()
			//var tmpdata Data
			//json.Unmarshal(tmp,&tmpdata)
			//h.roomID[c]=tmpdata.RoomId
			//fmt.Println("c.data.roomid is",c.data.RoomId)
			//fmt.Println("c.data",tmpdata)
			c.data.Ip = c.ws.RemoteAddr().String()
			c.data.Type = "handshake"
			c.data.UserList = user_list
			data_b, _ := json.Marshal(c.data)
			c.message_chan <- data_b
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				delete(h.roomID,c)
				close(c.message_chan)
			}
		case data := <-h.broadcast:
			for c := range h.connections {
				var tmp Data
				json.Unmarshal(data,&tmp)
				fmt.Println("tmp is ",tmp)
				fmt.Println("h[c] is",h.roomID[c])
				if tmp.RoomId==h.roomID[c]{
					select {
					case c.message_chan <- data:
					default:
						delete(h.connections, c)
						close(c.message_chan)
					}
				}

			}
		}
	}
}