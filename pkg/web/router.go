package web

import (
    "embed"
    "github.com/gin-gonic/gin"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/util"
    "github.com/yhy0/SuWen/pkg/web/controller"
    "github.com/yhy0/SuWen/pkg/web/middleware"
    "github.com/yhy0/logging"
    "html/template"
    "io/fs"
    "net/http"
    "strings"
    "time"
)

/**
   @author yhy
   @since 2024/5/21
   @desc //TODO
**/

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

func Run(port string) {
    gin.SetMode("release")
    router := gin.Default()
    
    // 创建一个新的模板引擎实例
    tmpl := template.New("")
    
    // 向模板引擎添加自定义函数
    tmpl.Funcs(template.FuncMap{
        "nl2br":          util.Nl2br,
        "splitString":    util.SplitString,
        "contains":       util.Contains,
        "pointer":        util.Pointer,
        "TimeSub":        util.TimeSub,
        "ParseMarkdown":  util.ParseMarkdown,
        "TruncateString": util.TruncateString,
    })
    
    router.Static("/db/results/", "./db/results/")
    
    // 静态资源加载
    router.StaticFS("/static", mustFS())
    
    // 使用自定义的模板引擎实例加载模板
    router.SetHTMLTemplate(template.Must(tmpl.ParseFS(templates, "templates/*")))
    
    visitCounter := middleware.NewVisitCounter()
    go visitCounter.UpdateDatabase()
    // 使用中间件
    router.Use(middleware.VisitCountMiddleware(visitCounter))
    
    router.GET("/", func(c *gin.Context) {
        c.Redirect(302, "/index")
    })
    
    router.GET("/index", func(c *gin.Context) {
        pressReleases := db.GetPressReleaseCategory()
        sourceNames := make(map[string]int)
        for _, pr := range pressReleases {
            if sourceNames[pr.SourceName] != 0 {
                sourceNames[pr.SourceName] += 1
            } else {
                sourceNames[pr.SourceName] = 1
            }
        }
        
        PieLabels, pieSeries := util.Sort(sourceNames)
        
        advisoryDate := db.GetAdvisoryDate()
        dateMap := make(map[string][]int)
        var chatLabels []string
        for _, date := range advisoryDate {
            _date := date.PublishedAt.Format("2006-01-02")
            if len(dateMap[_date]) != 0 {
                dateMap[_date][0] += 1
            } else {
                dateMap[_date] = []int{1, 0}
                chatLabels = append(chatLabels, _date)
            }
        }
        
        pocDate := db.GetPocDate()
        for _, date := range pocDate {
            _date := date.CommitDate.Format("2006-01-02")
            if len(dateMap[_date]) != 0 && dateMap[_date][1] != 0 {
                dateMap[_date][1] += 1
            } else if len(dateMap[_date]) == 0 { // 这种情况表明 cve 这天没有数据
                dateMap[_date] = []int{0, 1}
                chatLabels = append(chatLabels, _date)
            } else {
                dateMap[_date][1] = 1
            }
        }
        
        chatLabels = util.SortTimeSlice(util.UniqueStrings(chatLabels))
        
        var cveSeries []int
        var pocSeries []int
        var cveTotal, pocTotal int
        for _, label := range chatLabels {
            cveSeries = append(cveSeries, dateMap[label][0])
            cveTotal += dateMap[label][0]
            pocSeries = append(pocSeries, dateMap[label][1])
            pocTotal += dateMap[label][1]
        }
        
        pathMap := make(map[string]int)
        addrMap := make(map[string]int)
        ipMap := make(map[string]int)
        var pathSum int
        for _, pc := range db.GetPathCounts() {
            k := pc.Path
            if k != "/index" && k != "/about" && k != "/pr" && k != "/gad" && k != "/poc" && k != "/rss" {
                pathMap["other"] = pathMap["other"] + pc.Count
                pathSum = pathSum + pc.Count
            } else {
                pathMap[k] = pc.Count
                pathSum += pc.Count
            }
        }
        
        for _, ipc := range db.GetIPCounts() {
            if ipc.Address == "" {
                ipc.Address = "未知"
            }
            address := strings.Split(ipc.Address, " ")
            if len(address) == 2 {
                ipc.Address = address[0]
            }
            if _, exists := addrMap[ipc.Address]; !exists {
                addrMap[ipc.Address] = 1
            }
            addrMap[ipc.Address] += 1
            if _, exists := ipMap[ipc.Address]; !exists {
                ipMap[ipc.IP] = 1
            }
            ipMap[ipc.IP] += 1
        }
        
        PathLabels, PathSeries := util.Sort(pathMap)
        
        AddrLabels, AddrSeries := util.Sort(addrMap)
        
        c.HTML(http.StatusOK, "index.html", gin.H{
            "year":          time.Now().Year(),
            "version":       conf.Version,
            "active":        "home",
            "PieSeries":     pieSeries,
            "PieLabels":     PieLabels,
            "ChatLabels":    chatLabels,
            "CveSeries":     cveSeries,
            "PocSeries":     pocSeries,
            "cveTotal":      cveTotal,
            "pocTotal":      pocTotal,
            "PathLabels":    PathLabels,
            "PathSeries":    PathSeries,
            "AddrLabels":    AddrLabels,
            "AddrSeries":    AddrSeries,
            "LastCheckTime": conf.LastCheckTime,
            "pathSum":       pathSum,
            "ipCountSum":    len(ipMap),
            "ipAddrSum":     len(addrMap),
        })
    })
    
    router.GET("/about", func(c *gin.Context) {
        c.HTML(http.StatusOK, "about.html", gin.H{
            "year":    time.Now().Year(),
            "version": conf.Version,
            "active":  "about",
        })
    })
    
    // 设置路由
    router.GET("/rss", func(c *gin.Context) {
        rssBytes, err := controller.GenerateRSS()
        if err != nil {
            c.String(http.StatusInternalServerError, "Internal Server Error")
            return
        }
        c.Data(http.StatusOK, "application/xml", []byte(rssBytes))
    })
    router.GET("/pr", controller.GetPressRelease)
    router.POST("/pr", controller.GetPressRelease)
    
    router.GET("/gad", controller.GetGitHubAdvisoryDatabase)
    
    router.GET("/poc", controller.GetPoc)
    
    logging.Logger.Println("server start at port:", port)
    
    err := router.Run(":" + port)
    if err != nil {
        logging.Logger.Errorln(err)
        return
    }
}

func mustFS() http.FileSystem {
    sub, err := fs.Sub(static, "static")
    
    if err != nil {
        panic(err)
    }
    
    return http.FS(sub)
}
