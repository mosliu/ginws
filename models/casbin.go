package models

import (
    "github.com/jinzhu/gorm"
    "github.com/mosliu/ginws/db"
    "github.com/mosliu/ginws/middleware/casbin"
)

//定义一个 User 的结构体, 用来存放用户名和密码
type CasbinModel struct {
    gorm.Model //加入此行用于在数据库中创建记录的 mate 数据
    Ptype  string `gorm:"type:varchar(100);column:ptype" json:"ptype"`
    RoleName   string `gorm:"type:varchar(100);column:v0" json:"rolename"`
    Path   string `gorm:"type:varchar(100);column:v1" json:"path"`
    Method string `gorm:"type:varchar(100);column:v2" json:"method"`
    V3 string `gorm:"type:varchar(100);column:v3" json:"V3"`
    V4 string `gorm:"type:varchar(100);column:v4" json:"V4"`
    V5 string `gorm:"type:varchar(100);column:v5" json:"V5"`
}
func (CasbinModel) TableName() string {
    return "casbin_rule"
}


func init(){
    //AutoMigrate 只会添加缺失的字段，不会删除或者修改当前的字段
    db.AuthDB.AutoMigrate(&CasbinModel{})
}


func  (c *CasbinModel) AddCasbin(cm CasbinModel) bool {
    e := casbin.Enforcer

    return e.AddPolicy(cm.RoleName, cm.Path, cm.Method)

}
