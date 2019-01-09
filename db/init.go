package db
import (
    "github.com/json-iterator/go"
    "github.com/mosliu/ginws/logs"
    "github.com/sirupsen/logrus"
)

var log = logs.Log.WithFields(logrus.Fields{
    "pkg":"db",
})
var json = jsoniter.ConfigCompatibleWithStandardLibrary

