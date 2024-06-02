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

type Advisory struct {
    gorm.Model
    Id          int       `gorm:"primary_key" json:"id"`
    GhsaId      string    `json:"ghsa_id"`
    Summary     string    `json:"summary"`
    Description string    `json:"description"`
    Severity    string    `gorm:"index" json:"severity"`
    CVE         string    `gorm:"index" json:"cve"`
    Score       float64   `json:"score"`
    Ecosystem   string    `json:"ecosystem"`
    PublishedAt time.Time `json:"published_at"`
    GithubUrl   string    `json:"github_url"`
    Pushed      bool      `json:"pushed"`
}

func AddAdvisory(data *Advisory) int {
    globalDBTmp := GlobalDB.Model(&Advisory{})
    globalDBTmp.Create(&data)
    return data.Id
}

func GetAdvisoryInfo(ecosystem string) (count int64, advisories []*Advisory) {
    globalDBTmp := GlobalDB.Model(&Advisory{})
    if ecosystem != "" {
        globalDBTmp = globalDBTmp.Or("ecosystem = ?", ecosystem)
    }
    
    globalDBTmp.Count(&count).Order("published_at desc").Find(&advisories)
    
    return
}

func GetEcosystem() (advisories []*Advisory) {
    globalDBTmp := GlobalDB.Model(&Advisory{})
    globalDBTmp.Select("ecosystem").Find(&advisories)
    
    return
}

func GetAdvisoryDate() (advisories []*Advisory) {
    globalDBTmp := GlobalDB.Model(&Advisory{})
    globalDBTmp.Select("published_at").Find(&advisories)
    
    return
}

func SearchAdvisory(ghsaId string) bool {
    globalDBTmp := GlobalDB.Model(&Advisory{})
    var data Advisory
    globalDBTmp.Where("ghsa_id = ?", ghsaId).First(&data)
    if data.Id > 0 {
        return true
    }
    
    return false
}

func UpdateAdvisory(key string, data Advisory) {
    globalDBTmp := GlobalDB.Model(&Advisory{})
    globalDBTmp.Where("unique_key = ?", key).Updates(data)
}
