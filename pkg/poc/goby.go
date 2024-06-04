package poc

import (
    "encoding/json"
    "github.com/imroc/req/v3"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/util"
    "github.com/yhy0/logging"
    "strconv"
    "time"
)

/**
  @author: yhy
  @since: 2024/6/4
  @desc: 从 goby 官网获取 goby 新增的 poc 信息
**/

type gobyMsg struct {
    Data []struct {
        ReleaseVersion string `json:"release_version"`
        CreatedAt      string `json:"created_at"`
        CreatedAtWeek  string `json:"created_at_week"`
        PushData       []struct {
            Name             string `json:"name"`
            Description      string `json:"description"`
            Product          string `json:"product"`
            Impact           string `json:"impact"`
            Recommendation   string `json:"recommendation"`
            Tags             string `json:"tags"`
            NameEn           string `json:"name_en"`
            DescriptionEn    string `json:"description_en"`
            ProductEn        string `json:"product_en"`
            ImpactEn         string `json:"impact_en"`
            RecommendationEn string `json:"recommendation_en"`
            TagsEn           string `json:"tags_en"`
            CveId            string `json:"cve_id"`
            FofaQuery        string `json:"fofa_query"`
            AssetCount       int    `json:"asset_count"`
            Cvss             string `json:"cvss"`
            DemoGifUrl       string `json:"demo_gif_url"`
            ReleasedAt       []struct {
                Name   string `json:"name"`
                NameEn string `json:"name_en"`
            } `json:"released_at"`
        } `json:"push_data"`
    } `json:"data"`
    
    StatusCode int    `json:"statusCode"`
    Messages   string `json:"messages"`
}

func FindGoby() {
    c := req.C().
        SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36").
        SetTimeout(time.Duration(10) * time.Second).ImpersonateChrome()
    
    c.SetCommonRetryCount(3).
        SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second)
    
    resp, err := c.R().Get("https://gobysec.net/api/poc-push-list")
    if err != nil {
        logging.Logger.Error(err)
        return
    }
    
    if resp.StatusCode == 200 {
        var body gobyMsg
        err = resp.UnmarshalJson(&body)
        logging.Logger.Infoln("get goby pocs success.")
        
        for _, gobyPoc := range body.Data {
            if ok, t := util.IsToday(gobyPoc.CreatedAt); ok {
                for _, _poc := range gobyPoc.PushData {
                    if db.SearchPoc(_poc.Name) { // 只要新增的 poc
                        continue
                    }
                    var severity string
                    
                    cvss, _ := strconv.ParseFloat(_poc.Cvss, 64)
                    if cvss >= 8 || _poc.AssetCount > 3000 {
                        severity = "Critical"
                    } else {
                        severity = "High"
                    }
                    
                    poc := &db.Poc{
                        Description: _poc.Description,
                        PocName:     _poc.Name,
                        PocUrl:      "https://gobysec.net/updates",
                        Source:      "goby",
                        Severity:    severity,
                        CommitDate:  t,
                    }
                    db.AddPoc(poc)
                    pushPoc(poc)
                    data, _ := json.MarshalIndent(poc, "", "  ")
                    logging.Logger.Infoln("[+] Add poc", string(data))
                }
                
            }
        }
    }
    
}
