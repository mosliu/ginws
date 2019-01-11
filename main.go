package main

import (
    "github.com/fatih/color"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/mosliu/ginws/db"
    "github.com/mosliu/ginws/ginutils"
    "github.com/mosliu/ginws/middleware/jwt"
    "github.com/mosliu/ginws/webget"
    "github.com/mosliu/ginws/wsutils"
    "github.com/spf13/viper"
    "os"
    "os/signal"
    "sync"
    "syscall"

    //这个包用来实现一个 HTTP 的 web 框架
    "net/http"
)

//退出用 等待组
var waitGroupForExit sync.WaitGroup
//退出清理操作。
func doExit() {
    log.Infoln("Programme exiting....")
    //do sth
    log.Infoln("Programme exited.")
    os.Exit(0)
}

//解析退出信号
func parseSig(sigChan chan os.Signal) {
    waitGroupForExit.Add(1)
    for sig := range sigChan {
        switch sig {
        case syscall.SIGHUP, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
            log.Infoln("Receive signal:", sig, ",do exit ops")
            defer waitGroupForExit.Done()
            doExit()
        default:
            log.Infoln("Receive signal:", sig, " do nothing.")
        }
    }
}


//main 包必要一个 main 函数，作为起点
func main() {
    //信号程道
    sigChan := make(chan os.Signal)
    //监听指定的信号
    signal.Notify(sigChan)

    go parseSig(sigChan)

    jwt.InitDbAndCasbin()

    //func Default() *Engine
    //Default returns an Engine instance with the Logger and Recovery middleware already attached.
    //用来返回一个已经加载了Logger and Recovery中间件的引擎
    r := gin.Default()

    //the jwt middleware
    authMiddleware := jwt.AddMidd(viper.GetString("jwt.realm"), viper.GetString("jwt.key"), viper.GetString("jwt.tokenLookup"), viper.GetString("jwt.tokenHeadName"))

    //func (mw *GinJWTMiddleware) LoginHandler(c *gin.Context)
    //LoginHandler can be used by clients to get a jwt token. Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}. Reply will be of the form {"token": "TOKEN"}.
    //将 /login 交给 authMiddleware.LoginHandler 函数来处理
    r.POST(viper.GetString("jwt.loginPath"), authMiddleware.LoginHandler)
    //func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup
    //Group creates a new router group. You should add all the routes that have common middlwares or the same path prefix. For example, all the routes that use a common middlware for authorization could be grouped
    //创建一个组 auth
    auth := r.Group(viper.GetString("jwt.authPath"))
    //func (mw *GinJWTMiddleware) MiddlewareFunc() gin.HandlerFunc
    //MiddlewareFunc makes GinJWTMiddleware implement the Middleware interface.
    //auth 组中使用 MiddlewareFunc 中间件
    auth.Use(authMiddleware.MiddlewareFunc())
    {
        //如果是 /auth 组下的 /hello 就交给 helloHandler 来处理
        auth.GET(viper.GetString("jwt.testPath"), jwt.HelloHandler)
        //func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context)
        //RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh. Shall be put under an endpoint that is using the GinJWTMiddleware. Reply will be of the form {"token": "TOKEN"}.
        //如果是 /auth 组下的 /refresh_token 就交给 RefreshHandler 来处理
        auth.GET(viper.GetString("jwt.refreshPath"), authMiddleware.RefreshHandler)
        auth.GET("/ping", jwt.H_ping)
    }

    r.Run(":" + viper.GetString("server.port")) //在 0.0.0.0:配置端口 上启监听

    //退出等待
    waitGroupForExit.Wait()
}

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
    c := &wsutils.Client{SendChan: make(chan []byte, 256), WsConn: ws1}
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

func main2() {
    //test()
    log.Infof("%s starting", viper.GetString("name"))
    log.Warnln("file?")
    log.Debugln("color?")
    log.Debugln(color.BlueString("color!"))
    //r := gin.Default()
    go wsutils.CommonHub.Run()
    r := gin.New()
    r.Use(ginutils.Logger(log.Logger), gin.Recovery())
    r.Static("/assert", "./assert")
    r.GET("/ping2", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.GET("/ping", ping)
    r.GET("/wsutils", wsHandler)
    //bindAddress := "localhost:2303"
    bindAddress := viper.GetString("server.address") + ":" + viper.GetString("server.port")

    log.Infoln("Listening on ", bindAddress)

    webget.TransTKL("Aa")
    //crawlers.GetTrade()
    r.Run(bindAddress)
}

func main3() {
    //crawlers.GetTrade()
    db.DoSqlit3()
}
