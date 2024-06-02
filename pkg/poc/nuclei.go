package poc

import (
    "context"
    "encoding/json"
    "github.com/google/go-github/v62/github"
    "github.com/hashicorp/go-multierror"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/notice"
    "github.com/yhy0/SuWen/pkg/util"
    "github.com/yhy0/logging"
    "strings"
    "time"
)

/**
  @author: yhy
  @since: 2024/5/30
  @desc: //TODO
**/

func FindNucleiPR() {
    logging.Logger.Infoln("finding nuclei PR")
    var nucleiPR []*github.PullRequest
    
    ctx := context.Background()
    // 检查200个pr
    for page := 1; page < 2; page++ {
        prs, _, err := conf.GithubClient.PullRequests.List(ctx, "projectdiscovery", "nuclei-templates", &github.PullRequestListOptions{
            State:       "all",
            ListOptions: github.ListOptions{Page: page, PerPage: 100},
        })
        if err != nil {
            if len(prs) == 0 {
                logging.Logger.Warnf("list nuclei pr failed: %v", err)
                return
            } else {
                logging.Logger.Warnf("list nuclei pr failed: %v", err)
                continue
            }
        }
        nucleiPR = append(nucleiPR, prs...)
    }
    
    for _, pr := range nucleiPR {
        if today, _ := util.IsTodayBeijing(pr.GetCreatedAt().Time); today {
            files, _, err := conf.GithubClient.PullRequests.ListFiles(ctx, "projectdiscovery", "nuclei-templates", pr.GetNumber(), nil)
            if err != nil {
                logging.Logger.Warnf("get afrog commit failed: %v", err)
                continue
            }
            
            for _, file := range files {
                content := strings.ToLower(file.GetPatch())
                if file.GetStatus() == "added" && (strings.Contains(content, "severity: critical") || strings.Contains(content, "severity: high")) {
                    // 只要新增的 poc
                    if db.SearchPoc(file.GetFilename()) {
                        continue
                    }
                    loc, _ := time.LoadLocation("Asia/Shanghai")
                    var severity string
                    if strings.Contains(content, "severity: critical") {
                        severity = "Critical"
                    } else if strings.Contains(content, "severity: high") {
                        severity = "High"
                    }
                    poc := &db.Poc{
                        Description: pr.GetTitle(),
                        PocName:     file.GetFilename(),
                        PocUrl:      file.GetBlobURL(),
                        Source:      "nuclei-templates",
                        Severity:    severity,
                        CommitDate:  pr.GetCreatedAt().Time.In(loc),
                    }
                    data, _ := json.MarshalIndent(poc, "", "  ")
                    logging.Logger.Infoln("[+] Add poc", string(data))
                    db.AddPoc(poc)
                    pushPoc(poc)
                }
            }
        }
    }
    
}

func pushPoc(poc *db.Poc) {
    var pushErr *multierror.Error
    
    if err := notice.Text.PushMarkdown(poc.Description, notice.RenderPocMsg(poc)); err != nil {
        pushErr = multierror.Append(pushErr, err)
    }
    
    if err := notice.Raw.PushRaw(notice.NewRawPocMsgMessage(poc)); err != nil {
        pushErr = multierror.Append(pushErr, err)
    }
    
    if pushErr != nil {
        logging.Logger.Errorf("failed to push %s, %s", poc.PocName, pushErr.Error())
    }
}
