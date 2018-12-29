package ws
import (
    "github.com/mosliu/ginws/logs"
    "github.com/sirupsen/logrus"
)

var log = logs.Log.WithFields(logrus.Fields{
    "pkg":"ws",
})

var logf = logs.Log.WithFields(logrus.Fields{
    "pkg":"ws",
})