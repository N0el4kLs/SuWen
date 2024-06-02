package poc

import (
    "context"
    "encoding/json"
    "github.com/google/go-github/v62/github"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/db"
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

func FindAfrog() {
    logging.Logger.Infoln("finding Afrog commit")
    ctx := context.Background()
    // 只用检查检查 50 个 commit 就行了，应该不会一天更新超过这么多次
    afrogCommit, _, err := conf.GithubClient.Repositories.ListCommits(ctx, "zan8in", "afrog", &github.CommitsListOptions{
        Path:        "pocs/afrog-pocs/",
        ListOptions: github.ListOptions{Page: 1, PerPage: 50},
    })
    if err != nil {
        logging.Logger.Warnf("list afrog commit failed: %v", err)
        return
    }
    
    for _, commit := range afrogCommit {
        if today, _ := util.IsTodayBeijing(commit.Commit.Committer.GetDate().Time); today {
            _commit, _, err := conf.GithubClient.Repositories.GetCommit(ctx, "zan8in", "afrog", *commit.SHA, nil)
            if err != nil {
                logging.Logger.Warnf("get afrog commit failed: %v", err)
                continue
            }
            files := make(map[string]string)
            for _, file := range _commit.Files {
                content := strings.ToLower(file.GetPatch())
                if file.GetStatus() == "added" && (strings.Contains(content, "severity: critical") || strings.Contains(content, "severity: high")) {
                    if db.SearchPoc(file.GetFilename()) { // 只要新增的 poc
                        continue
                    }
                    files[file.GetFilename()] = file.GetBlobURL()
                    loc, _ := time.LoadLocation("Asia/Shanghai")
                    
                    var severity string
                    if strings.Contains(content, "severity: critical") {
                        severity = "Critical"
                    } else if strings.Contains(content, "severity: high") {
                        severity = "High"
                    }
                    
                    poc := &db.Poc{
                        Description: commit.Commit.GetMessage(),
                        PocName:     file.GetFilename(),
                        PocUrl:      file.GetBlobURL(),
                        Source:      "afrog",
                        Severity:    severity,
                        CommitDate:  commit.Commit.Committer.GetDate().Time.In(loc),
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
