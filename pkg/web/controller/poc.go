package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/util"
    "net/http"
    "strconv"
    "time"
)

/**
  @author: yhy
  @since: 2024/5/31
  @desc: //TODO
**/

func GetPoc(c *gin.Context) {
    startTime := time.Now()
    
    pageSize, _ := strconv.Atoi(c.Query("pageSize"))
    if pageSize == 0 {
        pageSize = 20
    }
    
    pageNum := 0
    page, _ := strconv.Atoi(c.Query("current"))
    if page == 0 {
        page = 1
    } else if page > 0 {
        pageNum = (page - 1) * pageSize
    }
    
    pocName := c.Query("pocName")
    
    totalCount, pocs := db.GetPocInfo(pageNum, pageSize, pocName)
    
    endTime := time.Now()
    
    // 计算两个时间点之间的差值
    duration := endTime.Sub(startTime)
    
    paginator := util.NewPaginator(c.Request, pageSize, totalCount)
    
    // 获取时间间隔的秒数
    seconds := duration.Seconds()
    
    c.HTML(http.StatusOK, "poc.html", gin.H{
        "year":       time.Now().Year(),
        "version":    conf.Version,
        "active":     "poc",
        "totalCount": totalCount,
        "seconds":    seconds,
        "pocs":       pocs,
        "paginator":  paginator,
        "pageSize":   pageSize,
    })
}
