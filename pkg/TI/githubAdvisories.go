package TI

import (
    "context"
    "encoding/json"
    "github.com/google/go-github/v62/github"
    "github.com/hashicorp/go-multierror"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/llm"
    "github.com/yhy0/SuWen/pkg/notice"
    "github.com/yhy0/SuWen/pkg/util"
    "github.com/yhy0/logging"
    "net/http"
    "strconv"
    "strings"
    "time"
)

/**
  @author: yhy
  @since: 2024/5/27
  @desc: https://github.com/advisories
**/

func RunGithubAdvisories() {
    FetchAdvisories()
    conf.LastCheckTime["Github Advisories"] = util.TimeNow()
    interval, err := time.ParseDuration(conf.GlobalConfig.WatchVulnAppConfig.Interval)
    if err != nil {
        logging.Logger.Errorln(err)
        return
    }
    
    if interval.Minutes() < 1 {
        interval, err = time.ParseDuration("1m")
    }
    
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for {
        logging.Logger.Infof("Github Advisories next checking at %s\n", time.Now().Add(interval).Format("2006-01-02 15:04:05"))
        
        select {
        case <-ticker.C:
            FetchAdvisories()
            conf.LastCheckTime["Github Advisories"] = util.TimeNow()
        }
    }
}

func FetchAdvisories() {
    tr := http.DefaultTransport.(*http.Transport).Clone()
    tr.Proxy = http.ProxyFromEnvironment
    githubClient := github.NewClient(&http.Client{
        Timeout:   time.Second * 10,
        Transport: tr,
    }).WithAuthToken(conf.GlobalConfig.GithubToken)
    
    opts := &github.ListGlobalSecurityAdvisoriesOptions{
        ListCursorOptions: github.ListCursorOptions{
            Page:    "1",
            PerPage: 100,
        },
    }
    
    var page = 1
    for {
        advisories, _, err := githubClient.SecurityAdvisories.ListGlobalSecurityAdvisories(context.Background(), opts)
        if err != nil {
            logging.Logger.Errorln(err)
            break
        }
        var cnt int
        for i, advisory := range advisories {
            // 只获取今天的
            ok, publishedAt := util.IsTodayBeijing(advisory.PublishedAt.Time)
            if !ok {
                continue
            }
            if *advisory.Severity != "high" && *advisory.Severity != "critical" {
                continue
            }
            // 将结构体转换为 JSON 字符串
            jsonData, _ := json.Marshal(advisory)
            if db.SearchAdvisory(advisory.GetGHSAID()) {
                logging.Logger.Infoln("[*]", string(jsonData))
                continue
            }
            logging.Logger.Infoln("[+]", string(jsonData))
            
            _, updatedAt := util.IsTodayBeijing(advisory.UpdatedAt.Time)
            _, githubReviewedAt := util.IsTodayBeijing(advisory.GithubReviewedAt.Time)
            advisory.PublishedAt.Time = publishedAt
            advisory.UpdatedAt.Time = updatedAt
            advisory.GithubReviewedAt.Time = githubReviewedAt
            cnt = i
            
            ecosystem := ""
            if advisory.Vulnerabilities != nil {
                v := advisory.Vulnerabilities[0]
                ecosystem = v.Package.GetEcosystem()
            }
            
            var score float64
            
            if advisory.GetCVSS().GetScore() != nil {
                score = *advisory.GetCVSS().GetScore()
            }
            
            description := advisory.GetDescription()
            
            _description := llm.ChatGPT(description)
            if _description != "" {
                description = _description
            }
            
            cve := &db.Advisory{
                GhsaId:      advisory.GetGHSAID(),
                Summary:     advisory.GetSummary(),
                Description: strings.ReplaceAll(description, "请帮我把以下内容翻译成中文，同时对输出的 markdown 语法和链接信息 前后空一格显示，并且只返回你处理后的内容，不要增加任何无关输出", ""),
                Severity:    advisory.GetSeverity(),
                CVE:         advisory.GetCVEID(),
                Score:       score,
                Ecosystem:   ecosystem,
                Pushed:      true,
                PublishedAt: advisory.PublishedAt.Time,
                GithubUrl:   advisory.GetHTMLURL(),
            }
            db.AddAdvisory(cve)
            pushCVE(cve)
        }
        if cnt == 99 {
            page += 1
            opts.ListCursorOptions.Page = strconv.Itoa(page)
        } else {
            break
        }
    }
}

func pushCVE(cve *db.Advisory) {
    var pushErr *multierror.Error
    
    if err := notice.Text.PushMarkdown(cve.Summary, notice.RenderCVEMsg(cve)); err != nil {
        pushErr = multierror.Append(pushErr, err)
    }
    
    if err := notice.Raw.PushRaw(notice.NewRawCVEMsgMessage(cve)); err != nil {
        pushErr = multierror.Append(pushErr, err)
    }
    
    if pushErr != nil {
        logging.Logger.Errorf("failed to push %s, %s", cve.CVE, pushErr.Error())
    }
}
