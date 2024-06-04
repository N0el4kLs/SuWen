package db

import (
    "gorm.io/gorm"
    "time"
)

/**
  @author: yhy
  @since: 2024/5/21
  @desc: //TODO
**/

type Poc struct {
    gorm.Model
    Id          int       `gorm:"primary_key" json:"id"`
    Description string    `json:"description"`
    PocName     string    `json:"poc_name"`
    PocUrl      string    `json:"poc_url"`
    Source      string    `gorm:"index" json:"source"`
    CommitDate  time.Time `gorm:"index" json:"commit_date"`
    Severity    string    `gorm:"index" json:"severity"`
    Pushed      bool      `json:"pushed"`
}

func AddPoc(data *Poc) int {
    globalDBTmp := GlobalDB.Model(&Poc{})
    globalDBTmp.Create(&data)
    return data.Id
}

func GetPocInfo(pageNum int, pageSize int, pocName string) (count int64, poc []*Poc) {
    globalDBTmp := GlobalDB.Model(&Poc{})
    
    if pocName != "" {
        globalDBTmp.Where("poc_name = ?", pocName)
    }
    
    globalDBTmp.Count(&count)
    
    if pageNum == 0 && pageSize == 0 {
        globalDBTmp.Order("commit_date desc").Find(&poc)
    } else {
        globalDBTmp.Offset(pageNum).Limit(pageSize).Order("commit_date desc").Find(&poc)
    }
    
    return
}

func SearchPoc(pocName string) bool {
    globalDBTmp := GlobalDB.Model(&Poc{})
    var data Poc
    globalDBTmp.Where("poc_name = ?", pocName).First(&data)
    if data.Id > 0 {
        return true
    }
    
    return false
}

func GetPocDate() (poc []*Poc) {
    globalDBTmp := GlobalDB.Model(&Poc{})
    globalDBTmp.Select("commit_date").Find(&poc)
    
    return
}
