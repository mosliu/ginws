package wsutils

import (
    "github.com/gorilla/websocket"
    "github.com/mosliu/ginws/webget"
)

type Connection struct {
    Id string
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
        ok,str :=webget.TransTKL(string(message))
        if ok{
            CommonHub.Broadcast <- []byte(str)
        }

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

