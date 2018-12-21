package logs

import (
    "runtime"
    "github.com/sirupsen/logrus"
    "strings"
    "fmt"
)

//
//type ContextHook struct {
//}
//
//func (hook ContextHook) Levels() []logrus.Level {
//    return logrus.AllLevels
//}
//func (hook ContextHook) Fire(entry *logrus.Entry) error {
//    if pc, file, line, ok := runtime.Caller(7); ok {
//        funcName := runtime.FuncForPC(pc).Name()
//        entry.Data["file"] = path.Base(file)
//        entry.Data["func"] = path.Base(funcName)
//        entry.Data["line"] = line
//    }
//    return nil
//}

type contextHook struct {
    Field  string
    Skip   int
    levels []logrus.Level
}
// NewContextHook use to make an hook
// 根据上面的推断, 我们递归深度可以设置到5即可.
func NewContextHook(levels ...logrus.Level) logrus.Hook {
    hook := contextHook{
        Field:  "source",
        Skip:   5,
        levels: levels,
    }
    if len(hook.levels) == 0 {
        hook.levels = logrus.AllLevels
    }
    return &hook
}
// Levels implement levels
func (hook contextHook) Levels() []logrus.Level {
    return logrus.AllLevels
}
// Fire implement fire
func (hook contextHook) Fire(entry *logrus.Entry) error {
    entry.Data[hook.Field] = findCaller(hook.Skip)
    return nil
}
// 对caller进行递归查询, 直到找到非logrus包产生的第一个调用.
// 因为filename我获取到了上层目录名, 因此所有logrus包的调用的文件名都是 logrus/...
// 因此通过排除logrus开头的文件名, 就可以排除所有logrus包的自己的函数调用
func findCaller(skip int) string {
    file := ""
    line := 0
    for i := 0; i < 10; i++ {
        file, line = getCaller(skip + i,2)
        if !strings.HasPrefix(file, "logrus") {
            file, line = getCaller(skip + i,1)
            break
        }
    }

    //return color.MagentaString("%s:%d", file, line)
    return fmt.Sprintf("%s:%d", file, line)
}
// 这里其实可以获取函数名称的: fnName := runtime.FuncForPC(pc).Name()
// 但是我觉得有 文件名和行号就够定位问题, 因此忽略了caller返回的第一个值:pc
// 在标准库log里面我们可以选择记录文件的全路径或者文件名, 但是在使用过程成并发最合适的,
// 因为文件的全路径往往很长, 而文件名在多个包中往往有重复, 因此这里选择多取一层, 取到文件所在的上层目录那层.
// remain :返回的string保留的路径层数，例如logrus/entry 为2层  entry为1层，因此在判断logrus包时需要设定为2
func getCaller(skip int,remain int) (string, int) {
    _, file, line, ok := runtime.Caller(skip)
    //fmt.Println(file)
    //fmt.Println(line)
    if !ok {
        return "", 0
    }
    n := 0
    for i := len(file) - 1; i > 0; i-- {
        if file[i] == '/' {
            n++
            if n >= remain {
                file = file[i+1:]
                break
            }
        }
    }
    return file, line
}