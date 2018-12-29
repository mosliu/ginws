package ginutils

import (
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
)

// 你可以创建很多instance
//Log to stdout.
var Log = logrus.New()

// Log to File.
var LogF = logrus.New()

type Logs struct {
    log  *logrus.Logger
    logF *logrus.Logger
}

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
