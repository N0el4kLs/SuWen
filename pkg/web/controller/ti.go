package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/util"
    "net/http"
    "time"
)

/**
   @author yhy
   @since 2024/5/22
   @desc //TODO
**/

func GetPressRelease(c *gin.Context) {
    startTime := time.Now()
    formSourceName := c.Query("sourceName")
    formTags := c.PostFormArray("form-tags")
    formSeverity := c.PostFormArray("form-severity")
    
    search := make(map[string][]string)
    if formSourceName != "" {
        if formSourceName == "all" {
            formSourceName = ""
        }
        search["source_name"] = []string{formSourceName}
    }
    if len(formTags) > 0 {
        search["tags"] = formTags
    }
    if len(formSeverity) > 0 {
        search["severity"] = formSeverity
    }
    
    totalCount, pressReleases := db.GetPressRelease(search)
    tags := make(map[string]int)
    sourceNames := make(map[string]int)
    severity := make(map[string]int)
    
    _pressReleases := db.GetPressReleaseCategory()
    for _, pr := range _pressReleases {
        _tags := util.SplitString(pr.Tags, " ,")
        for _, _tag := range _tags {
            if _tag != "" {
                if tags[_tag] != 0 {
                    tags[_tag] += 1
                } else {
                    tags[_tag] = 1
                }
            }
        }
        if sourceNames[pr.SourceName] != 0 {
            sourceNames[pr.SourceName] += 1
        } else {
            sourceNames[pr.SourceName] = 1
        }
        if severity[pr.Severity] != 0 {
            severity[pr.Severity] += 1
        } else {
            severity[pr.Severity] = 1
        }
    }
    
    sourceNames["all"] = len(_pressReleases)
    
    endTime := time.Now()
    
    // 计算两个时间点之间的差值
    duration := endTime.Sub(startTime)
    
    // 获取时间间隔的秒数
    seconds := duration.Seconds()
    
    if formSourceName == "" {
        formSourceName = "all"
    }
    c.HTML(http.StatusOK, "pr.html", gin.H{
        "year":          time.Now().Year(),
        "version":       conf.Version,
        "active":        "ti",
        "totalCount":    totalCount,
        "seconds":       seconds,
        "PressReleases": pressReleases,
        "tags":          tags,
        "sourceNames":   sourceNames,
        "sourceName":    formSourceName,
        "severity":      severity,
    })
}

func GetGitHubAdvisoryDatabase(c *gin.Context) {
    startTime := time.Now()
    ecosystem := c.Query("ecosystem")
    if ecosystem == "all" {
        ecosystem = ""
    }
    
    totalCount, advisories := db.GetAdvisoryInfo(ecosystem)
    ecosystems := make(map[string]int)
    _advisories := db.GetEcosystem()
    for _, advisory := range _advisories {
        if ecosystems[advisory.Ecosystem] != 0 {
            ecosystems[advisory.Ecosystem] += 1
        } else if advisory.Ecosystem != "" {
            ecosystems[advisory.Ecosystem] = 1
        }
    }
    ecosystems["all"] = len(_advisories)
    
    endTime := time.Now()
    
    // 计算两个时间点之间的差值
    duration := endTime.Sub(startTime)
    
    // 获取时间间隔的秒数
    seconds := duration.Seconds()
    if ecosystem == "" {
        ecosystem = "all"
    }
    c.HTML(http.StatusOK, "advisory.html", gin.H{
        "year":       time.Now().Year(),
        "version":    conf.Version,
        "active":     "ti",
        "totalCount": totalCount,
        "seconds":    seconds,
        "advisories": advisories,
        "ecosystems": ecosystems,
        "ecosystem":  ecosystem,
    })
}
