package wsutils
import (
    "github.com/mosliu/ginws/logs"
    "github.com/sirupsen/logrus"
)

var log = logs.Log.WithFields(logrus.Fields{
    "pkg":"wsutils",
})

var logf = logs.Log.WithFields(logrus.Fields{
    "pkg":"wsutils",
})