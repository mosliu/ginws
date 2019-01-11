package casbin

import (
    "github.com/mosliu/ginws/logs"
    "github.com/sirupsen/logrus"
)

// 你可以创建很多instance
//Log to stdout.
var log = logs.Log.WithFields(logrus.Fields{
    "pkg":"casbin",
})