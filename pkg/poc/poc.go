package poc

import (
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/util"
    "github.com/yhy0/logging"
    "time"
)

/**
  @author: yhy
  @since: 2024/5/30
  @desc: //TODO
**/

func Run() {
    FindNucleiPR()
    conf.LastCheckTime["Nuclei-Templates"] = util.TimeNow()
    
    FindAfrog()
    conf.LastCheckTime["Afrog"] = util.TimeNow()
    
    FindGoby()
    conf.LastCheckTime["Goby"] = util.TimeNow()
    
    interval, err := time.ParseDuration(conf.GlobalConfig.WatchVulnAppConfig.Interval)
    if err != nil {
        logging.Logger.Errorln(err)
        return
    }
    
    if interval.Minutes() < 1 {
        interval, err = time.ParseDuration("15m")
    }
    
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for {
        logging.Logger.Infof("poc next checking at %s\n", time.Now().Add(interval).Format("2006-01-02 15:04:05"))
        
        select {
        case <-ticker.C:
            FindNucleiPR()
            conf.LastCheckTime["Nuclei-Templates"] = util.TimeNow()
            
            FindAfrog()
            conf.LastCheckTime["Afrog"] = util.TimeNow()
            
            FindGoby()
            conf.LastCheckTime["Goby"] = util.TimeNow()
            
            // hour := time.Now().In(loc).Hour()
            // if hour >= 0 && hour < 7 {
            //     // we must sleep in this time
            //     logging.Logger.Infof("sleeping..")
            //     continue
            // }
        }
    }
}
