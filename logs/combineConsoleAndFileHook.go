package logs

import (
    "github.com/sirupsen/logrus"
)

type combineHook struct {
}

// Levels implement levels
func (hook combineHook) Levels() []logrus.Level {
    return logrus.AllLevels
}

// Fire implement fire
func (hook combineHook) Fire(entry *logrus.Entry) error {

    logF.WithFields(entry.Data)
    //logF.Level = entry.Level
    switch entry.Level {
    case logrus.PanicLevel:
        logF.Panic(entry.Message)
        //return hook.Writer.Crit(line)
    case logrus.FatalLevel:
        logF.Fatal(entry.Message)
        //return hook.Writer.Crit(line)
    case logrus.ErrorLevel:
        logF.Error(entry.Message)
        //return hook.Writer.Err(line)
    case logrus.WarnLevel:
        logF.Warn(entry.Message)
        //return hook.Writer.Warning(line)
    case logrus.InfoLevel:
        logF.Info(entry.Message)
        //return hook.Writer.Info(line)
    case logrus.DebugLevel:
        logF.Debug(entry.Message)
        //return hook.Writer.Debug(line)
    default:
        return nil
    }
    return nil
}

func NewCombineConsoleAndFile() logrus.Hook {
    hook := combineHook{}
    return hook
}
