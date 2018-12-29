package logs

import (
    "github.com/lestrrat-go/file-rotatelogs"
    "github.com/mosliu/ginws/utils"
    "github.com/rifflock/lfshook"
    "github.com/shiena/ansicolor"
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
    "os"
    "path/filepath"
    "time"
)

// 你可以创建很多instance
//Log to stdout.
var Log = logrus.New()

// Log to File.
var logF = logrus.New()

func init() {
    if !utils.GetConfigLoadStatus() {
        utils.LoadConfig()
    }
    logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
    logrus.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
    //levelConsole := "debug"
    levelConsole := viper.GetString("logs.level.console")
    //levelFile := "info"
    levelFile := viper.GetString("logs.level.file")

    logPath := viper.GetString("logs.file.path")
    logFilename := viper.GetString("logs.file.name")
    maxcount := viper.GetInt("logs.file.maxcount")
    logFileSuffix := viper.GetString("logs.file.suffix")
    initConsoleLogger(levelConsole)
    initFileLogger(levelFile, logPath, logFilename,logFileSuffix, uint(maxcount))
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

    Log.SetOutput(colorWriter)

    Log.AddHook(NewContextHook())
    Log.AddHook(NewCombineConsoleAndFile())

    Log.Info("Console Logger Initialized.")
}

func initFileLogger(level string, logPath string, logFilename string,logFileSuffix string, maxcount uint) {
    //先检查目录
    exist, err := utils.PathExists(logPath)
    if err != nil {
        logrus.Fatalln("Log Path Error", err)
        return
    }

    //无则创建
    if !exist {
        logrus.Infof("Log path %v is not exist,create it.\r\n")
        // 创建文件夹
        err := os.MkdirAll(logPath, os.ModePerm)
        if err != nil {
            logrus.Fatalln("Log Path create Error.", err)
        }
    }

    //init LogF
    //logFilename := "./logrus.Log"
    logfilefullname := filepath.Join(logPath, logFilename)
    //file, err := os.OpenFile(logfilefullname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    //if err != nil {
    //    logrus.Errorln("Failed to Log to file, using default stderr")
    //    return
    //}

    writer, err := rotatelogs.New(
        logfilefullname+".%Y%m%d."+logFileSuffix,
        rotatelogs.WithLinkName(logfilefullname+"."+logFileSuffix),
        // WithMaxAge和WithRotationCount二者只能设置一个，
        // WithMaxAge设置文件清理前的最长保存时间，
        // WithRotationCount设置文件清理前最多保存的个数。
        //rotatelogs.WithMaxAge(time.Hour*24),
        rotatelogs.WithRotationCount(maxcount),
        // WithRotationTime设置日志分割的时间，这里设置为一天分割一次
        rotatelogs.WithRotationTime(24*time.Hour),
    )

    if err != nil {
        logrus.Errorln("config local file system for logger error: %v", err)
    }

    //设定level
    lvlFile, err := logrus.ParseLevel(level)
    if err != nil {
        Log.Fatal(err)
        logF.SetLevel(logrus.InfoLevel)
    } else {
        logF.SetLevel(lvlFile)
    }

    lfsHook := lfshook.NewHook(lfshook.WriterMap{
        logrus.DebugLevel: writer,
        logrus.InfoLevel:  writer,
        logrus.WarnLevel:  writer,
        logrus.ErrorLevel: writer,
        logrus.FatalLevel: writer,
        logrus.PanicLevel: writer,
    }, &logrus.JSONFormatter{})

    logF.AddHook(lfsHook)
    //logF.SetOutput(file)

    //logF.SetFormatter(&logrus.JSONFormatter{})


    Log.Info("File Logger Initialized.")
    //LogF.Info("Logger Component Initialized.")
}
