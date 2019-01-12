package db

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/spf13/viper"
)

//定义一个内部全局的 db 指针用来进行认证，数据校验
//var AuthDB = &gorm.DB{}
var AuthDB *gorm.DB
//定义一个内部全局的 db 指针用来进行认证，数据校验
var MainDB *gorm.DB

func init() {
    log.Infof("AuthDB:%x", AuthDB)
    AuthDB = InitDb("db.authdb")
    log.Infof("AuthDB:%x", AuthDB)
    MainDB = InitDb("db")
}

// 初始化数据库，使用配置中confPrefix前缀的配置
func InitDb(confPrefix string) *gorm.DB {
    dbType := viper.GetString(confPrefix + ".db_type")
    var dbToOpen *gorm.DB
    var err error
    switch dbType {
    case "sqlite3":
        sqlitePath := viper.GetString(confPrefix + ".sqlite_path")
        dbToOpen, err = gorm.Open(dbType, sqlitePath)
    case "mysql":
        viper.SetDefault(confPrefix+".db_charset", "utf8")
        mysqlConnInfo := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
            viper.GetString(confPrefix+".db_user"), viper.GetString(confPrefix+".db_password"),
            viper.GetString(confPrefix+".db_host"), viper.GetString(confPrefix+".db_port"),
            viper.GetString(confPrefix+".db_name"), viper.GetString(confPrefix+".db_charset"),
        )
        dbToOpen, err = gorm.Open(dbType, mysqlConnInfo)
    case "postgres":
        viper.SetDefault(confPrefix+".db_sslmode", "disable")
        pgConnInfo := fmt.Sprintf(
            "host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
            viper.GetString(confPrefix+".db_host"), viper.GetString(confPrefix+".db_port"),
            viper.GetString(confPrefix+".db_user"), viper.GetString(confPrefix+".db_name"),
            viper.GetString(confPrefix+".db_password"), viper.GetString(confPrefix+".db_sslmode"),
        )
        dbToOpen, err = gorm.Open(dbType, pgConnInfo)
    case "none":
        log.Infoln("No Support DB Provided")
    default:
        log.Errorln("No Support DB Provided")
    }

    //sqlite_info := viper.GetString("db.sqlitePath")
    //opendb, err := gorm.Open("sqlite3", sqlite_info)
    if err != nil {
        //如果出错，就直接打印出错信息，并且退出
        log.Panicln(err)
    } else {
        fmt.Println("Successfully connected!") //如果没有出错，就打印成功连接的信息
    }
    //AuthDB.LogMode(true)
    viper.SetDefault(confPrefix+".log_mode", false)
    dbToOpen.LogMode(viper.GetBool(confPrefix + ".log_mode"))
    dbToOpen.SetLogger(log)
    //用于设置最大打开的连接数，默认值为0表示不限制。
    viper.SetDefault(confPrefix+".max_conn", 200)
    dbToOpen.DB().SetMaxOpenConns(viper.GetInt(confPrefix + ".max_conn"))
    //用于设置闲置的连接数。
    viper.SetDefault(confPrefix+".idle_con", 100)
    dbToOpen.DB().SetMaxOpenConns(viper.GetInt(confPrefix + ".idle_con"))
    dbToOpen.DB().Ping()
    return dbToOpen
}

func CloseDbs() {
    AuthDB.Close()
    MainDB.Close()
}
