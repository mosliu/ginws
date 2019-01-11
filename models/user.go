package models

import (
    "github.com/jinzhu/gorm"
    "github.com/mosliu/ginws/db"
    "github.com/spf13/viper"
)

//定义一个 User 的结构体, 用来存放用户名和密码
type User struct {
    gorm.Model //加入此行用于在数据库中创建记录的 mate 数据
    UserID   string `gorm:"type:varchar(30);UNIQUE;unique_index;not null" form:"username" json:"username" binding:"required"`
    Password string `gorm:"size:255" form:"password" json:"password" binding:"required"`
}

func init(){
    //AutoMigrate 只会添加缺失的字段，不会删除或者修改当前的字段
    db.AuthDB.AutoMigrate(&User{})
}

func defaultUser(){
    user := User{}
    //如果 db 中没有这条记录，就创建，如果有就忽略掉
    viper.SetDefault("rbac.admin_name","root")
    viper.SetDefault("rbac.admin_pass","password")
    //db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
    notExists := db.AuthDB.Where(&User{UserID: viper.GetString("rbac.admin_name")}).First(&user).RecordNotFound()
    if notExists {
        user := User{UserID: viper.GetString("rbac.admin_name"), Password: viper.GetString("rbac.admin_pass")}
        db.AuthDB.Create(&user)
    }

}