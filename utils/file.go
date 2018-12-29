package utils

import (
    "errors"
    "os"
)

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
    dir, err := os.Stat(path)

    if err == nil {
        if dir.IsDir(){
            return true, nil
        }else{
            return false, errors.New("not a path")
        }
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

