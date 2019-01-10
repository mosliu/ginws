package ginutils

import (
    "github.com/gin-gonic/gin"
    "github.com/mosliu/ginws/logs"
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
)

// 你可以创建很多instance
//Log to stdout.
var log = logs.Log.WithFields(logrus.Fields{
    "pkg":"ginutils",
})

func init() {
    setGinMode()
}

func setGinMode(){
    mode:= viper.GetString("gin.mode")
    switch mode{
    case "release":
        gin.SetMode(gin.ReleaseMode)
    case "debug":
        gin.SetMode(gin.DebugMode)
    default:
        gin.SetMode(gin.DebugMode)
    }

}
