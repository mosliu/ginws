package utils

import "github.com/sirupsen/logrus"

func FatalError(err error,logto *logrus.Entry){
    if err != nil{
        logto.WithField("error",err).Fatal("Error Occurs")
    }

}

func CheckError(err error,logto *logrus.Entry){
    if err != nil{
        logto.WithField("error",err).Warn("Error Occurs")
    }
}