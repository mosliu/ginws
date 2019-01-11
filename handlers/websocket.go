package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/mosliu/ginws/wsutils"
    "net/http"
)



var upGrader = &websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}


func WsHandler(ctx *gin.Context) {
    ws1, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
    if err != nil {
        return
    }
    c := &wsutils.Client{SendChan: make(chan []byte, 256), WsConn: ws1}
    wsutils.CommonHub.Register <- c
    defer func() { wsutils.CommonHub.Unregister <- c }()
    go c.Writer()
    c.Reader()
}

//webSocket 请求 ping 返回 pong
func WsPing(ctx *gin.Context) {
    // 升级 get 请求为 webSocket 协议
    ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
    if err != nil {
        return
    }
    defer ws.Close()
    for {
        // 读取 wsutils 中的数据
        mt, message, err := ws.ReadMessage()
        if err != nil {
            break
        }
        if string(message) == "ping" {
            message = []byte("pong")
        }
        // 写入 wsutils 数据
        err = ws.WriteMessage(mt, message)
        if err != nil {
            break
        }
    }
}
