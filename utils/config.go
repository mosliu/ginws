package utils

import (
    "fmt"
    "github.com/fsnotify/fsnotify"
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
    "os"
)

var isConfigLoaded = false

func GetConfigLoadStatus() bool {
    return isConfigLoaded
}
func LoadConfig() {
    if isConfigLoaded {
        return
    }

    //viper.SetEnvPrefix(cmdRoot)
    viper.AutomaticEnv()
    //replacer := strings.NewReplacer(".", "_")
    //viper.SetEnvKeyReplacer(replacer)

    //load server
    configFileName := "server"
    viper.SetConfigName(configFileName) //  设置配置文件名 (不带后缀)
    viper.AddConfigPath("./")
    viper.AddConfigPath("./configs")

    err := viper.ReadInConfig()
    if err != nil {
        fmt.Println(fmt.Errorf("Fatal error when reading %s config file:%s", configFileName, err))
        os.Exit(1)
    }

    exists,isdir,err := PathFileExists("./configs/personal.toml")
    if err ==nil {
        if exists && (!isdir){
            pfile,_:=os.Open("./configs/personal.toml")
            viper.MergeConfig(pfile)
        }
    }

    //监视配置文件，重新读取配置数据
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        logrus.Warnln("Config file changed:", e.Name)
        //fmt.Println("Config file changed:", e.Name)
    })

    name := viper.GetString("tbk.apkey")
    fmt.Println("apk:", name)
    //consoleLevel := viper.GetString("logs.console.level")
    //fmt.Println("logs.console.level:", consoleLevel)
    //
    //fullstate := viper.GetString("statetransfer.timeout.fullstate")
    //fmt.Println("statetransfer.timeout.fullstate:", fullstate)
    //
    //abcdValue := viper.GetString("peer.abcd")
    //fmt.Println("peer.abcd:", abcdValue)

    isConfigLoaded = true
}
