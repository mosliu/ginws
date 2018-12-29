package main

import (
    "github.com/fatih/color"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/mosliu/ginws/ginutils"
    "github.com/mosliu/ginws/wsutils"
    "github.com/spf13/viper"
    "net/http"
)

var upGrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func wsHandler(ctx *gin.Context) {
    ws1, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
    if err != nil {
        return
    }
    c := &wsutils.Connection{SendChan: make(chan []byte, 256), WsConn: ws1}
    wsutils.CommonHub.Register <- c
    defer func() { wsutils.CommonHub.Unregister <- c }()
    go c.Writer()
    c.Reader()
}

//webSocket 请求 ping 返回 pong
func ping(ctx *gin.Context) {
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

func main() {
    //test()
    log.Infof("%s starting",viper.GetString("name"))
    log.Warnln("file?")
    log.Debugln("color?")
    log.Debugln(color.BlueString("color!"))
    //r := gin.Default()
    go wsutils.CommonHub.Run()
    r := gin.New()
    r.Use(ginutils.Logger(log.Logger), gin.Recovery())
    r.Static("/assert","./assert")
    r.GET("/ping2", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.GET("/ping", ping)
    r.GET("/wsutils", wsHandler)
    //bindAddress := "localhost:2303"
    bindAddress := viper.GetString("server.address") + ":" + viper.GetString("server.port")
    r.Run(bindAddress)
    log.Infoln("Listening on ",bindAddress)
}
