package utils

import (
    "fmt"
    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
    "os"
)


func LoadConfig() {
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

    //监视配置文件，重新读取配置数据
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        fmt.Println("Config file changed:", e.Name)
    })

    //environment := viper.GetBool("security.enabled")
    //fmt.Println("security.enabled:", environment)
    //
    //fullstate := viper.GetString("statetransfer.timeout.fullstate")
    //fmt.Println("statetransfer.timeout.fullstate:", fullstate)
    //
    //abcdValue := viper.GetString("peer.abcd")
    //fmt.Println("peer.abcd:", abcdValue)
}
