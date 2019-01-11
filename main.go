package main

import (
    "github.com/gin-gonic/gin"
    "github.com/mosliu/ginws/db"
    "github.com/mosliu/ginws/ginutils"
    "github.com/mosliu/ginws/handlers"
    "github.com/mosliu/ginws/middleware/casbin"
    "github.com/mosliu/ginws/middleware/jwt"
    "github.com/spf13/viper"
    "os"
    "os/signal"
    "sync"
    "syscall"
)

//退出用 等待组
var waitGroupForExit sync.WaitGroup
//退出清理操作。
func doExit() {
    log.Infoln("Programme exiting....")
    //do sth
    //close db connections
    db.CloseDbs()
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

// 配置 启动 gin server
func startGinServer() {
    //初始化casbin 大写 非自动
    casbin.Init()

    //用来返回一个已经加载了Logger and Recovery中间件的引擎
    r := gin.New()
    r.Use(ginutils.Logger(log.Logger), gin.Recovery())

    //使用 jwt 中间件
    authMiddleware := jwt.GetJwtMidlleware()

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
        auth.GET(viper.GetString("jwt.testPath"), handlers.HelloHandler)
        //func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context)
        //RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh. Shall be put under an endpoint that is using the GinJWTMiddleware. Reply will be of the form {"token": "TOKEN"}.
        //如果是 /auth 组下的 /refresh_token 就交给 RefreshHandler 来处理
        auth.GET(viper.GetString("jwt.refreshPath"), authMiddleware.RefreshHandler)
        auth.GET("/ping", handlers.H_ping)
    }

    //非认证
    r.GET("/wsutils", handlers.WsHandler)
    r.GET("/wsping", handlers.WsPing)

    //静态资源
    r.Static("/assert", "./assert")

    bindAddress := viper.GetString("server.address") + ":" + viper.GetString("server.port")

    log.Infoln("Listening on ", bindAddress)

    r.Run(bindAddress) //在 0.0.0.0:配置端口 上启监听

}

//main 包必要一个 main 函数，作为起点
func main() {
    //信号程道
    sigChan := make(chan os.Signal)
    //监听指定的信号
    signal.Notify(sigChan)
    //信号处理
    go parseSig(sigChan)

    go startGinServer()

    //退出等待
    waitGroupForExit.Wait()
}
