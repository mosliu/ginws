package db

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DgtleTradInfo struct{
    gorm.Model
    Title string
    UserId string
    Price uint
    Url string

}

func DoSqlit3(){
    db, err := gorm.Open("sqlite3", "./gorm.db")
    if err!=nil{
        log.Error(err)
    }
    defer db.Close()
    item1 := DgtleTradInfo{
        Title:"aaa",
        UserId:"dadd",
        Price:90,
        Url:"http://www.baidu.com",

    }

    db.AutoMigrate(&DgtleTradInfo{})

    if db.NewRecord(item1){
        db.Create(&item1)
    }


}