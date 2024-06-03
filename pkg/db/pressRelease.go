package db

import (
    "gorm.io/gorm"
)

/**
  @author: yhy
  @since: 2024/5/21
  @desc: //TODO
**/

// PressRelease 厂商 PR 漏洞情报信息表
type PressRelease struct {
    gorm.Model
    Id           int    `gorm:"primary_key" json:"id"`
    UniqueKey    string `gorm:"index" json:"unique_key"`
    Title        string `json:"title"`
    Description  string `json:"description"`
    Severity     string `gorm:"index" json:"severity"`
    CVE          string `gorm:"index" json:"cve"`
    Disclosure   string `json:"disclosure"`
    Solutions    string `json:"solutions"`
    GithubSearch string `json:"github_search"`
    References   string `json:"references"`
    Pushed       bool   `gorm:"index" json:"pushed"`
    Tags         string `json:"tags"`
    From         string `json:"from"`
    Source       string `json:"source"`
    SourceName   string `json:"source_name"`
    Color        string `json:"color"`
}

func AddPressRelease(data *PressRelease) int {
    globalDBTmp := GlobalDB.Model(&PressRelease{})
    globalDBTmp.Create(&data)
    return data.Id
}

func GetPressRelease(query map[string][]string) (count int64, PressReleases []*PressRelease) {
    if query != nil {
        searchPressRelease(query).Count(&count).Order("updated_at desc").Find(&PressReleases)
    } else {
        GlobalDB.Model(&PressRelease{}).Count(&count).Order("updated_at desc").Find(&PressReleases)
    }
    
    return
}

func GetPressReleaseCategory() (PressReleases []*PressRelease) {
    globalDBTmp := GlobalDB.Model(&PressRelease{})
    globalDBTmp.Select("tags, source_name, severity").Find(&PressReleases)
    return
}

func SearchPressReleaseByKey(key string) (data *PressRelease) {
    globalDBTmp := GlobalDB.Model(&PressRelease{})
    globalDBTmp.Where("unique_key = ?", key).First(&data)
    return
}

func SearchPressReleaseByCve(cve string) (data []*PressRelease) {
    globalDBTmp := GlobalDB.Model(&PressRelease{})
    globalDBTmp.Where("cve = ? and pushed = ?", cve, true).Find(&data)
    return
}

func UpdatePressReleaseByKey(key, links string) {
    globalDBTmp := GlobalDB.Model(&PressRelease{})
    globalDBTmp.Where("unique_key = ?", key).Update("github_search", links)
}

func UpdatePressReleasePushed(key string) {
    globalDBTmp := GlobalDB.Model(&PressRelease{})
    globalDBTmp.Where("unique_key = ?", key).Update("pushed", true)
}

func UpdatePressRelease(key string, data PressRelease) {
    globalDBTmp := GlobalDB.Model(&PressRelease{})
    globalDBTmp.Where("unique_key = ?", key).Updates(data)
}

func GetPressReleaseTotal() int {
    var count int64
    globalDBTmp := GlobalDB.Model(&PressRelease{})
    globalDBTmp.Count(&count)
    return int(count)
}

func searchPressRelease(query map[string][]string) *gorm.DB {
    globalDBTmp := GlobalDB.Model(&PressRelease{})
    if query["tags"] != nil {
        for i, tag := range query["tags"] {
            if i > 0 {
                globalDBTmp = globalDBTmp.Or("tags LIKE ?", "%"+tag+"%")
            } else {
                globalDBTmp = globalDBTmp.Where("tags LIKE ?", "%"+tag+"%")
            }
        }
    }
    
    if query["source_name"] != nil {
        for i, source := range query["source_name"] {
            if i > 0 {
                globalDBTmp = globalDBTmp.Or("source_name = ?", source)
            } else {
                globalDBTmp = globalDBTmp.Where("source_name = ?", source)
            }
        }
    }
    
    if query["title"] != nil {
        for i, title := range query["title"] {
            if i > 0 {
                globalDBTmp = globalDBTmp.Or("title LIKE ?", "%"+title+"%")
            } else {
                globalDBTmp = globalDBTmp.Where("title LIKE ?", "%"+title+"%")
            }
        }
    }
    
    if query["key"] != nil {
        globalDBTmp = globalDBTmp.Where("unique_key = ?", query["key"][0])
    }
    
    if query["severity"] != nil {
        for i, severity := range query["severity"] {
            if i > 0 {
                globalDBTmp = globalDBTmp.Or("severity = ?", severity)
            } else {
                globalDBTmp = globalDBTmp.Where("severity = ?", severity)
            }
        }
    }
    
    if query["cve"] != nil {
        for i, cve := range query["cve"] {
            if i > 0 {
                globalDBTmp = globalDBTmp.Or("cve LIKE ?", "%"+cve+"%")
            } else {
                globalDBTmp = globalDBTmp.Where("cve LIKE ?", "%"+cve+"%")
            }
        }
    }
    
    return globalDBTmp
}
