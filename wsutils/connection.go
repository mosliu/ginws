package ws

import (
    "github.com/gorilla/websocket"
)

type Connection struct {
    // websocket 连接器
    WsConn *websocket.Conn

    // 发送信息的缓冲 channel
    SendChan chan []byte
}

func (c *Connection) Reader() {
    for {
        _, message, err := c.WsConn.ReadMessage()
        if err != nil {
            break
        }
        CommonHub.Broadcast <- message
    }
    c.WsConn.Close()
}

func (c *Connection) Writer() {
    for message := range c.SendChan {
        err := c.WsConn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            break
        }
    }
    c.WsConn.Close()
}

