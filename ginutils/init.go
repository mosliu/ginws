package logs

import (
    "github.com/mosliu/ginws/utils"
    "github.com/shiena/ansicolor"
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
    "io"
    "os"
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
    if !utils.GetConfigLoadStatus(){
       utils.LoadConfig()
    }
    logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
    logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
    //levelConsole := "debug"
    levelConsole := viper.GetString("logs.level.console")
    //levelFile := "info"
    levelFile := viper.GetString("logs.level.file")
    initConsoleLogger(levelConsole)
    initFileLogger(levelFile)
    Assemble()
}

func initConsoleLogger(level string) {
    lvlConsole, err := logrus.ParseLevel(level)
    if err != nil {
        Log.Fatal(err)
        Log.SetLevel(logrus.DebugLevel)
    } else {
        //Log.Info("setLevel to :",lvlConsole)
        Log.SetLevel(lvlConsole)
        //Log.Debug("set Level to :",lvlConsole)
    }

    // force colors on for TextFormatter
    Log.Formatter = &logrus.TextFormatter{ForceColors: true}
    // then wrap the Log output with it
    // 用于解决windows的terminal中彩色不正确的问题
    colorWriter := ansicolor.NewAnsiColorWriter(os.Stdout)
    logrus.SetOutput(colorWriter)
}

func initFileLogger(level string) {
    //init LogF
    fileLocation := "./logrus.Log"
    file, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err == nil {
        //设定双途径输出
        colorWriter := ansicolor.NewAnsiColorWriter(os.Stdout)
        LogF.SetOutput(io.MultiWriter(file, colorWriter))
        //LogF.Out = file
    } else {
        Log.Warn("Failed to Log to file, using default stderr")
    }

    lvlFile, err := logrus.ParseLevel(level)
    if err != nil {
        Log.Fatal(err)
        LogF.SetLevel(logrus.InfoLevel)
    } else {
        LogF.SetLevel(lvlFile)
    }
    Log.Info("Logger Component Initialized.")
    //LogF.Info("Logger Component Initialized.")
}

//assemble hooks
func Assemble() {
    //logrus.AddHook(NewContextHook())
    Log.AddHook(NewContextHook())
    LogF.AddHook(NewContextHook())
}
