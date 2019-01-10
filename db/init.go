package db
import (
    "fmt"
    "github.com/jinzhu/gorm"
    "github.com/json-iterator/go"
    "github.com/mosliu/ginws/logs"
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
)

var log = logs.Log.WithFields(logrus.Fields{
    "pkg":"db",
})
var json = jsoniter.ConfigCompatibleWithStandardLibrary


//定义一个内部全局的 db 指针用来进行认证，数据校验
var AuthDB *gorm.DB

func init(){
    sqlite_info := viper.GetString("db.sqlite_path")
    opendb ,err := gorm.Open("sqlite3", sqlite_info)
    if err != nil {
        panic(err) //如果出错，就直接打印出错信息，并且退出
    } else {

        fmt.Println("Successfully connected!") //如果没有出错，就打印成功连接的信息
        AuthDB = opendb                        //连接成功的情况下将认证的数据库进行赋值
    }
    AuthDB.LogMode(true)
    AuthDB.SetLogger(log)
    //用于设置最大打开的连接数，默认值为0表示不限制。
    AuthDB.DB().SetMaxOpenConns(200)
    //用于设置闲置的连接数。
    AuthDB.DB().SetMaxIdleConns(100)
    AuthDB.DB().Ping()

}