package main

import (
	"github.com/mosliu/ginws/logs"
	"github.com/mosliu/ginws/utils"
	"github.com/sirupsen/logrus"
)

var log = logs.Log.WithFields(logrus.Fields{
	"pkg": "main",
})
//var logf = logs.LogF.WithFields(logrus.Fields{
//	"pkg": "main",
//})
func init() {
    utils.LoadConfig()
}