package controller

import (
    "fmt"
    "github.com/gorilla/feeds"
    "github.com/yhy0/SuWen/pkg/db"
    "time"
)

/**
   @author yhy
   @since 2024/6/1
   @desc 支持 rss 订阅
**/

func GenerateRSS() (string, error) {
    now := time.Now()
    feed := &feeds.Feed{
        Title:       "素问",
        Link:        &feeds.Link{Href: "https://su.fireline.fun/"},
        Description: "漏洞信息推送",
        Author:      &feeds.Author{Name: "yhy", Email: "yhysec@qq.com"},
        Created:     now,
    }
    
    _, pocs := db.GetPocInfo(0, 0)
    for _, poc := range pocs {
        feed.Items = append(feed.Items, &feeds.Item{
            Title:       fmt.Sprintf("%s 又有新 POC 了", poc.Source),
            Link:        &feeds.Link{Href: "https://su.fireline.fun/poc"},
            Description: poc.Description,
            Created:     poc.CreatedAt,
        })
    }
    
    _, prs := db.GetPressRelease(nil)
    for _, pr := range prs {
        feed.Items = append(feed.Items, &feeds.Item{
            Title:       fmt.Sprintf("[%s] %s", pr.Source, pr.Title),
            Link:        &feeds.Link{Href: "https://su.fireline.fun/pr"},
            Description: pr.Description,
            Created:     pr.CreatedAt,
        })
    }
    
    _, advisories := db.GetAdvisoryInfo("")
    for _, advisory := range advisories {
        feed.Items = append(feed.Items, &feeds.Item{
            Title:       "又有新 CVE 了",
            Link:        &feeds.Link{Href: "https://su.fireline.fun/gad"},
            Description: advisory.Description,
            Created:     advisory.CreatedAt,
        })
    }
    
    return feed.ToRss()
}
