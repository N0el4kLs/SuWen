package controller

import (
    "fmt"
    "github.com/gorilla/feeds"
    "github.com/yhy0/SuWen/pkg/db"
    "sort"
)

/**
   @author yhy
   @since 2024/6/1
   @desc 支持 rss 订阅
**/

func GenerateRSS() (string, error) {
    feed := &feeds.Feed{
        Title:       "素问",
        Link:        &feeds.Link{Href: "https://su.fireline.fun/"},
        Description: "漏洞信息推送",
        Author:      &feeds.Author{Name: "yhy", Email: "yhysec@qq.com"},
    }
    
    var items []*feeds.Item
    
    _, pocs := db.GetPocInfo(0, 0, "")
    for _, poc := range pocs {
        items = append(items, &feeds.Item{
            Title:       fmt.Sprintf("%s 又有新 POC 了", poc.Source),
            Link:        &feeds.Link{Href: fmt.Sprintf("https://su.fireline.fun/poc?pocName=%s", poc.PocName)},
            Description: poc.Description,
            Created:     poc.CreatedAt,
        })
    }
    
    _, prs := db.GetPressRelease(nil)
    for _, pr := range prs {
        items = append(items, &feeds.Item{
            Title:       fmt.Sprintf("[%s] %s", pr.Source, pr.Title),
            Link:        &feeds.Link{Href: fmt.Sprintf("https://su.fireline.fun/pr?key=%s", pr.UniqueKey)},
            Description: pr.Description,
            Created:     pr.CreatedAt,
        })
    }
    
    _, advisories := db.GetAdvisoryInfo("", "")
    for _, advisory := range advisories {
        items = append(items, &feeds.Item{
            Title:       fmt.Sprintf("又有新 CVE 了 [%s]", advisory.CVE),
            Link:        &feeds.Link{Href: fmt.Sprintf("https://su.fireline.fun/gad?key=%s", advisory.GhsaId)},
            Description: advisory.Description,
            Created:     advisory.CreatedAt,
        })
    }
    
    // 按Created字段排序，最新的时间在前
    sort.Slice(items, func(i, j int) bool {
        return items[i].Created.After(items[j].Created)
    })
    
    feed.Items = items
    feed.Created = items[0].Created
    return feed.ToRss()
}
