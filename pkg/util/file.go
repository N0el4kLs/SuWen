package util

import (
    "os"
)

/**
  @author: yhy
  @since: 2024/5/21
  @desc: //TODO
**/

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
    _, err := os.Stat(path) // os.Stat获取文件信息
    if err != nil {
        if os.IsExist(err) {
            return true
        }
        return false
    }
    return true
}
