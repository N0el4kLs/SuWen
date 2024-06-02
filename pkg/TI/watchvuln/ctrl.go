package watchvuln

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/yhy0/SuWen/pkg/TI/watchvuln/grab"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/notice"
    "github.com/yhy0/SuWen/pkg/util"
    "github.com/yhy0/logging"
    "net/http"
    "regexp"
    "strings"
    "sync"
    "time"
    
    "github.com/google/go-github/v62/github"
    "github.com/hashicorp/go-multierror"
    "github.com/pkg/errors"
    "golang.org/x/sync/errgroup"
)

const (
    InitPageLimit   = 3
    UpdatePageLimit = 1
)

type WatchVulnApp struct {
    config       *conf.WatchVulnAppConfig
    githubClient *github.Client
    grabbers     []grab.Grabber
    prs          []*github.PullRequest
}

func NewApp(config *conf.WatchVulnAppConfig) (*WatchVulnApp, error) {
    var grabs []grab.Grabber
    for _, part := range strings.Split(config.Sources, ",") {
        part = strings.ToLower(strings.TrimSpace(part))
        switch part {
        case "avd":
            grabs = append(grabs, grab.NewAVDCrawler())
        case "nox", "ti":
            grabs = append(grabs, grab.NewTiCrawler())
        case "oscs":
            grabs = append(grabs, grab.NewOSCSCrawler())
        case "seebug":
            grabs = append(grabs, grab.NewSeebugCrawler())
        case "threatbook":
            grabs = append(grabs, grab.NewThreatBookCrawler())
        case "struts2", "structs2":
            grabs = append(grabs, grab.NewStruts2Crawler())
        case "kev":
            grabs = append(grabs, grab.NewKEVCrawler())
        default:
            return nil, fmt.Errorf("invalid grab source %s", part)
        }
    }
    
    tr := http.DefaultTransport.(*http.Transport).Clone()
    tr.Proxy = http.ProxyFromEnvironment
    githubClient := github.NewClient(&http.Client{
        Timeout:   time.Second * 10,
        Transport: tr,
    })
    
    return &WatchVulnApp{
        config:       config,
        githubClient: githubClient,
        grabbers:     grabs,
    }, nil
}

func (w *WatchVulnApp) Run() error {
    ctx := context.Background()
    defer ctx.Done()
    if w.config.DiffMode {
        logging.Logger.Info("running in diff mode, skip init vuln database")
        w.collectAndPush(ctx)
        logging.Logger.Info("diff finished")
        return nil
    }
    
    logging.Logger.Infof("initialize local database..")
    success, fail := w.initData(ctx)
    w.grabbers = success
    localCount := db.GetPressReleaseTotal()
    
    logging.Logger.Infof("system init finished, local database has %d vulns", localCount)
    if !w.config.NoStartMessage {
        providers := make([]*grab.Provider, 0, 10)
        failed := make([]*grab.Provider, 0, 10)
        for _, p := range w.grabbers {
            providers = append(providers, p.ProviderInfo())
        }
        for _, p := range fail {
            failed = append(failed, p.ProviderInfo())
        }
        msg := &notice.InitialMessage{
            VulnCount:      localCount,
            Interval:       w.config.Interval,
            Provider:       providers,
            FailedProvider: failed,
        }
        // 将结构体转换为 JSON 字符串
        data, err := json.MarshalIndent(msg, "", "  ")
        if err != nil {
            logging.Logger.Errorln("Error marshalling to JSON:", err)
        }
        
        logging.Logger.Infoln("WatchVuln 初始化完成", string(data))
        
        if err := notice.Text.PushMarkdown("WatchVuln 初始化完成", notice.RenderInitialMsg(msg)); err != nil {
            return err
        }
        if err := notice.Raw.PushRaw(notice.NewRawInitialMessage(msg)); err != nil {
            return err
        }
    }
    
    logging.Logger.Infof("ticking every %s", w.config.Interval)
    
    defer func() {
        msg := "注意: WatchVuln 进程退出"
        if err := notice.Text.PushText(msg); err != nil {
            logging.Logger.Error(err)
        }
        if err := notice.Raw.PushRaw(notice.NewRawTextMessage(msg)); err != nil {
            logging.Logger.Error(err)
        }
        time.Sleep(time.Second)
    }()
    
    interval, err := time.ParseDuration(conf.GlobalConfig.WatchVulnAppConfig.Interval)
    if err != nil {
        return err
    }
    
    if interval.Minutes() < 1 {
        interval, err = time.ParseDuration("1m")
    }
    
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for {
        w.prs = nil
        logging.Logger.Infof("watchvuln next checking at %s\n", time.Now().Add(interval).Format("2006-01-02 15:04:05"))
        
        select {
        case <-ticker.C:
            loc, _ := time.LoadLocation("Asia/Shanghai")
            hour := time.Now().In(loc).Hour()
            if hour >= 0 && hour < 7 {
                // we must sleep in this time
                logging.Logger.Infof("sleeping..")
                continue
            }
            w.collectAndPush(ctx)
        }
    }
}

func (w *WatchVulnApp) collectAndPush(ctx context.Context) {
    vulns, err := w.collectUpdate(ctx)
    if err != nil {
        logging.Logger.Errorf("failed to get updates, %s", err)
    }
    logging.Logger.Infof("found %d new vulns in this ticking", len(vulns))
    for _, v := range vulns {
        if w.config.NoFilter || v.Creator.IsValuable(v) {
            dbVuln := db.SearchPressReleaseByKey(v.UniqueKey)
            
            if dbVuln.Id <= 0 {
                logging.Logger.Errorf("%s from db not found %v", v.UniqueKey, dbVuln)
                continue
            }
            
            if dbVuln.Pushed {
                logging.Logger.Infof("%s has been pushed, skipped", v)
                continue
            }
            if v.CVE != "" && w.config.EnableCVEFilter {
                // 同一个 cve 已经有其它源推送过了
                others := db.SearchPressReleaseByCve(v.CVE)
                if len(others) == 0 {
                    logging.Logger.Errorf("%s from db not found %v", v.UniqueKey, others)
                    continue
                }
                
                ids := make([]string, 0, len(others))
                for _, o := range others {
                    ids = append(ids, o.UniqueKey)
                }
                logging.Logger.Infof("found new cve but other source has already pushed, others: %v", ids)
                continue
            }
            
            // find cve pr in nuclei repo
            if v.CVE != "" && !w.config.NoGithubSearch {
                links, err := w.FindGithubPoc(ctx, v.CVE)
                if err != nil {
                    logging.Logger.Warn(err)
                }
                logging.Logger.Infof("%s found %d links from github, %v", v.CVE, len(links), links)
                if len(links) != 0 {
                    v.GithubSearch = grab.MergeUniqueString(v.GithubSearch, links)
                    db.UpdatePressReleaseByKey(v.UniqueKey, strings.Join(v.GithubSearch, " ,"))
                }
            }
            logging.Logger.Infof("Pushing %s", v)
            // retry 3 times
            for i := 0; i < 3; i++ {
                if err := w.pushVuln(v); err == nil {
                    // 如果两种推送都成功，才标记为已推送
                    db.UpdatePressReleasePushed(v.UniqueKey)
                    break
                } else {
                    logging.Logger.Errorf("failed to push %s, %s", v.UniqueKey, err)
                }
                logging.Logger.Infof("retry to push %s after 30s", v.UniqueKey)
                time.Sleep(time.Second * 30)
            }
        } else {
            logging.Logger.Infof("skipped %s as not valuable", v)
        }
    }
}

func (w *WatchVulnApp) pushVuln(vul *grab.VulnInfo) error {
    var pushErr *multierror.Error
    
    if err := notice.Text.PushMarkdown(vul.Title, notice.RenderVulnInfo(vul)); err != nil {
        pushErr = multierror.Append(pushErr, err)
    }
    
    if err := notice.Raw.PushRaw(notice.NewRawVulnInfoMessage(vul)); err != nil {
        pushErr = multierror.Append(pushErr, err)
    }
    
    return pushErr.ErrorOrNil()
}

func (w *WatchVulnApp) initData(ctx context.Context) ([]grab.Grabber, []grab.Grabber) {
    var eg errgroup.Group
    eg.SetLimit(len(w.grabbers))
    var success []grab.Grabber
    var fail []grab.Grabber
    for _, grabber := range w.grabbers {
        gb := grabber
        eg.Go(func() error {
            source := gb.ProviderInfo()
            logging.Logger.Infof("start to init data from %s", source.Name)
            initVulns, err := gb.GetUpdate(ctx, InitPageLimit)
            if err != nil {
                fail = append(fail, gb)
                logging.Logger.Errorln(err)
                return errors.Wrap(err, source.Name)
            }
            
            for _, data := range initVulns {
                if _, err = w.createOrUpdate(source, data); err != nil {
                    fail = append(fail, gb)
                    return errors.Wrap(errors.Wrap(err, data.String()), source.Name)
                }
            }
            success = append(success, gb)
            return nil
        })
    }
    err := eg.Wait()
    if err != nil {
        logging.Logger.Error(errors.Wrap(err, "init data"))
    }
    return success, fail
}

func (w *WatchVulnApp) collectUpdate(ctx context.Context) ([]*grab.VulnInfo, error) {
    var eg errgroup.Group
    eg.SetLimit(len(w.grabbers))
    
    var mu sync.Mutex
    var newVulns []*grab.VulnInfo
    
    for _, grabber := range w.grabbers {
        gb := grabber
        eg.Go(func() error {
            source := gb.ProviderInfo()
            dataChan, err := gb.GetUpdate(ctx, UpdatePageLimit)
            if err != nil {
                return errors.Wrap(err, gb.ProviderInfo().Name)
            }
            
            conf.LastCheckTime[source.DisplayName] = util.TimeNow()
            
            hasNewVuln := false
            
            for _, data := range dataChan {
                isNewVuln, err := w.createOrUpdate(source, data)
                if err != nil {
                    return errors.Wrap(err, gb.ProviderInfo().Name)
                }
                if isNewVuln {
                    logging.Logger.Infof("found new vuln: %s", data)
                    mu.Lock()
                    newVulns = append(newVulns, data)
                    mu.Unlock()
                    hasNewVuln = true
                }
            }
            
            // 如果一整页漏洞都是旧的，说明没有更新，不必再继续下一页了
            if !hasNewVuln {
                return nil
            }
            return nil
        })
    }
    err := eg.Wait()
    return newVulns, err
}

func (w *WatchVulnApp) createOrUpdate(source *grab.Provider, data *grab.VulnInfo) (bool, error) {
    if !util.CurrentlyYear(data.Disclosure) || !util.MonthlyCalculation(data.Disclosure) { // 不是今年的或者相差 3 个月以上的就不要了
        return false, nil
    }
    
    pr := db.SearchPressReleaseByKey(data.UniqueKey)
    var color string
    switch source.Name {
    case "aliyun-avd":
        color = "silver"
    case "qianxin-ti":
        color = "azure"
    case "oscs":
        color = "indigo"
    case "threatbook":
        color = "red"
    case "seebug":
        color = "black"
    case "Struts2":
        color = "blue"
    case "KEV":
        color = "cyan"
    }
    
    var tags []string
    for _, tag := range data.Tags {
        if tag == "POC公开" {
            tag = "POC已公开"
        }
        if tag == "EXP公开" {
            tag = "EXP已公开"
        }
        tags = append(tags, tag)
    }
    data.Tags = tags
    // not exist
    if pr.Id <= 0 {
        data.Reason = append(data.Reason, grab.ReasonNewCreated)
        id := db.AddPressRelease(&db.PressRelease{
            UniqueKey:   data.UniqueKey,
            Title:       data.Title,
            Description: data.Description,
            Severity:    string(data.Severity),
            CVE:         data.CVE,
            Disclosure:  data.Disclosure,
            Solutions:   data.Solutions,
            References:  strings.Join(data.References, " ,"),
            Pushed:      false,
            Tags:        strings.Join(data.Tags, " ,"),
            From:        data.From,
            Source:      source.Name,
            SourceName:  source.DisplayName,
            Color:       color,
        })
        if id > 0 {
            logging.Logger.Infof("vuln %s(%s) created from %s %s", data.Title, data.UniqueKey, source.Name, data.From)
        } else {
            logging.Logger.Errorf("vuln %s(%s) not created from %s %s", data.Title, data.UniqueKey, source.Name, data.From)
        }
        return true, nil
    }
    
    // 如果一个漏洞之前是低危，后来改成了严重，这种可能也需要推送, 走一下高价值的判断逻辑
    asNewVuln := false
    if string(data.Severity) != pr.Severity {
        logging.Logger.Infof("%s from %s change severity from %s to %s", data.Title, data.From, pr.Severity, data.Severity)
        data.Reason = append(data.Reason, fmt.Sprintf("%s: %s => %s", grab.ReasonSeverityUpdated, pr.Severity, data.Severity))
        asNewVuln = true
    }
    for _, newTag := range data.Tags {
        found := false
        for _, dbTag := range strings.Split(pr.Tags, " ,") {
            if newTag == dbTag {
                found = true
                break
            }
        }
        // tag 有更新
        if !found {
            logging.Logger.Infof("%s from %s add new tag %s", data.Title, data.From, newTag)
            data.Reason = append(data.Reason, fmt.Sprintf("%s: %v => %v", grab.ReasonTagUpdated, pr.Tags, data.Tags))
            asNewVuln = true
            break
        }
    }
    
    if asNewVuln && data.Creator.IsValuable(data) && util.CurrentlyYear(data.Disclosure) { // 有很多老洞会更新导致最新的大都是老洞的信息，影响查看，这里判断一下是否有价值更新和披露时间
        // update
        db.UpdatePressRelease(data.UniqueKey, db.PressRelease{
            UniqueKey:   data.UniqueKey,
            Title:       data.Title,
            Description: data.Description,
            Severity:    string(data.Severity),
            CVE:         data.CVE,
            Disclosure:  data.Disclosure,
            Solutions:   data.Solutions,
            References:  strings.Join(data.References, " ,"),
            Pushed:      false,
            Tags:        strings.Join(data.Tags, " ,"),
            From:        data.From,
            Source:      source.Name,
            Color:       color,
        })
        
        logging.Logger.Debugf("vuln updated from %s %s", data.UniqueKey, source.Name)
    }
    
    return asNewVuln, nil
}

func (w *WatchVulnApp) FindGithubPoc(ctx context.Context, cveId string) ([]string, error) {
    var eg errgroup.Group
    var results []string
    var mu sync.Mutex
    
    eg.Go(func() error {
        links, err := w.findGithubRepo(ctx, cveId)
        if err != nil {
            return errors.Wrap(err, "find github repo")
        }
        mu.Lock()
        defer mu.Unlock()
        results = append(results, links...)
        return nil
    })
    eg.Go(func() error {
        links, err := w.findNucleiPR(ctx, cveId)
        if err != nil {
            return errors.Wrap(err, "find nuclei PR")
        }
        mu.Lock()
        defer mu.Unlock()
        results = append(results, links...)
        return nil
    })
    err := eg.Wait()
    return results, err
}

func (w *WatchVulnApp) findGithubRepo(ctx context.Context, cveId string) ([]string, error) {
    logging.Logger.Infof("finding github repo of %s", cveId)
    re, err := regexp.Compile(fmt.Sprintf(`(?i)(\b|_|/)%s(\b|_|/)`, cveId))
    if err != nil {
        return nil, err
    }
    lastYear := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
    query := fmt.Sprintf(`language:Python language:JavaScript language:C language:C++ language:Java language:PHP language:Ruby language:Rust language:C# created:>%s %s`, lastYear, cveId)
    result, _, err := w.githubClient.Search.Repositories(ctx, query, &github.SearchOptions{
        ListOptions: github.ListOptions{Page: 1, PerPage: 100},
    })
    if err != nil {
        return nil, err
    }
    var links []string
    for _, repo := range result.Repositories {
        if re.MatchString(repo.GetHTMLURL()) {
            links = append(links, repo.GetHTMLURL())
        }
    }
    return links, nil
}

func (w *WatchVulnApp) findNucleiPR(ctx context.Context, cveId string) ([]string, error) {
    logging.Logger.Infof("finding nuclei PR of %s", cveId)
    if w.prs == nil {
        // 检查200个pr
        for page := 1; page < 2; page++ {
            prs, _, err := w.githubClient.PullRequests.List(ctx, "projectdiscovery", "nuclei-templates", &github.PullRequestListOptions{
                State:       "all",
                ListOptions: github.ListOptions{Page: page, PerPage: 100},
            })
            if err != nil {
                if len(w.prs) == 0 {
                    return nil, err
                } else {
                    logging.Logger.Warnf("list nuclei pr failed: %v", err)
                    continue
                }
            }
            w.prs = append(w.prs, prs...)
        }
    }
    
    var links []string
    re, err := regexp.Compile(fmt.Sprintf(`(?i)(\b|_|/)%s(\b|_|/)`, cveId))
    if err != nil {
        return nil, err
    }
    for _, pr := range w.prs {
        if re.MatchString(pr.GetTitle()) || re.MatchString(pr.GetBody()) {
            links = append(links, pr.GetHTMLURL())
        }
    }
    return links, nil
}
