package utils

import (
    "errors"
    "os"
)

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
    dir, err := os.Stat(path)

    if err == nil {
        if dir.IsDir() {
            return true, nil
        } else {
            return false, errors.New("not a path")
        }
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

// 判断文件夹是否存在
// 返回是否存在，是否目录，错误
func PathFileExists(path string) (bool, bool, error) {
    dir, err := os.Stat(path)

    if err == nil {
        if dir.IsDir() {
            return true, true, nil
        } else {
            return true, false, nil
        }
    }
    if os.IsNotExist(err) {
        return false, false, nil
    }
    return false, false, err
}
