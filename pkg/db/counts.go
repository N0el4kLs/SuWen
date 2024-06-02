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

func AddOrUpdatePathCounts(path string, data *PathCounts) int {
    var count int
    var pc PathCounts
    GlobalDB.Model(&PathCounts{}).Where("path = ?", path).First(&pc)
    if pc.Id > 0 {
        if data.Count == 1 && pc.Count != 1 { // 这种应该是服务重启了，所以需要加 1
            data.Count = data.Count + pc.Count
            count = data.Count
        }
        GlobalDB.Model(&PathCounts{}).Where("path = ?", path).Updates(data)
    } else {
        GlobalDB.Model(&PathCounts{}).Create(data)
    }
    
    return count
}

func AddOrUpdateIPCounts(path, ip string, data *IPCounts) int {
    var count int
    var ic IPCounts
    GlobalDB.Model(&IPCounts{}).Where("path = ? and ip = ?", path, ip).First(&ic)
    if ic.Id > 0 {
        if data.Count == 1 && ic.Count != 1 { // 这种应该是服务重启了，所以需要加 1,然后将数据库中的数字返回
            data.Count = data.Count + ic.Count
            count = data.Count
        }
        GlobalDB.Model(&IPCounts{}).Where("path = ? and ip = ?", path, ip).Updates(data)
    } else {
        GlobalDB.Model(&IPCounts{}).Create(data)
    }
    
    return count
}
