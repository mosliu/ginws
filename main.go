package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/mosliu/ginws/ginutils"
    "github.com/mosliu/ginws/utils"
    "github.com/spf13/viper"
    "net/http"
)

var upGrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

//webSocket 请求 ping 返回 pong
func ping(c *gin.Context) {
    // 升级 get 请求为 webSocket 协议
    ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    defer ws.Close()
    for {
        // 读取 ws 中的数据
        mt, message, err := ws.ReadMessage()
        if err != nil {
            break
        }
        if string(message) == "ping" {
            message = []byte("pong")
        }
        // 写入 ws 数据
        err = ws.WriteMessage(mt, message)
        if err != nil {
            break
        }
    }
}

func main() {
    //test()
    utils.LoadConfig()
    log.Debugln(viper.GetString("name"))
    logf.Warnln("file?")
    log.Debugln(  "color?")
    //r := gin.Default()


    r := gin.New()
    r.Use(ginutils.Logger(log.Logger), gin.Recovery())
    r.GET("/ping2", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.GET("/ping", ping)
    //bindAddress := "localhost:2303"
    bindAddress := viper.GetString("server.address")+":"+viper.GetString("server.port")
    r.Run(bindAddress)
}
