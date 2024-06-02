package db

import (
    "gorm.io/gorm"
)

/**
   @author yhy
   @since 2024/6/2
   @desc //TODO
**/

type PathCounts struct {
    gorm.Model
    Id    int    `gorm:"primary_key" json:"id"`
    Path  string `gorm:"index" json:"path"`
    Count int    `json:"count"`
}

type IPCounts struct {
    gorm.Model
    Id      int    `gorm:"primary_key" json:"id"`
    IP      string `gorm:"index" json:"ip"`
    Path    string `gorm:"index" json:"path"`
    Count   int    `json:"count"`
    Address string `json:"address"`
}

func GetPathCounts() (pathCounts []*PathCounts) {
    GlobalDB.Model(&PathCounts{}).Order("count desc").Find(&pathCounts)
    return
}

func GetIPCounts() (ipCounts []*IPCounts) {
    GlobalDB.Model(&IPCounts{}).Order("count desc").Find(&ipCounts)
    return
}

func AddOrUpdatePathCounts(path string, data PathCounts) {
    globalDBTmp := GlobalDB.Model(&PathCounts{})
    globalDBTmp.Where("path = ?", path).FirstOrCreate(&data)
}

func AddOrUpdateIPCounts(path, ip string, data IPCounts) {
    globalDBTmp := GlobalDB.Model(&IPCounts{})
    globalDBTmp.Where("path = ? and ip = ?", path, ip).FirstOrCreate(&data)
}
