package utils

import "github.com/sirupsen/logrus"

func FatalError(err error,logto *logrus.Entry) error{
    if err != nil{
        logto.WithField("error",err).Fatal("Error Occurs")
    }
    return err

}

func CheckError(err error,logto *logrus.Entry) error{
    if err != nil{
        logto.WithField("error",err).Warn("Error Occurs")
    }
    return err
}