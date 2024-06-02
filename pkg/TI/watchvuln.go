package TI

import (
    "github.com/pkg/errors"
    "github.com/yhy0/SuWen/pkg/TI/watchvuln"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/logging"
)

/**
  @author: yhy
  @since: 2024/5/21
  @desc: //TODO
**/

func RunPR() error {
    logging.Logger.Infof("config: INTERVAL=%s, NO_FILTER=%v, NO_START_MESSAGE=%v, NO_GITHUB_SEARCH=%v, ENABLE_CVE_FILTER=%v",
        conf.GlobalConfig.WatchVulnAppConfig.Interval, conf.GlobalConfig.WatchVulnAppConfig.NoFilter, conf.GlobalConfig.WatchVulnAppConfig.NoStartMessage, conf.GlobalConfig.WatchVulnAppConfig.NoGithubSearch, conf.GlobalConfig.WatchVulnAppConfig.EnableCVEFilter)
    
    config := &conf.GlobalConfig.WatchVulnAppConfig
    
    app, err := watchvuln.NewApp(config)
    if err != nil {
        return errors.Wrap(err, "failed to create app")
    }
    
    if err = app.Run(); err != nil {
        return errors.Wrap(err, "failed to run app")
    }
    return nil
}
