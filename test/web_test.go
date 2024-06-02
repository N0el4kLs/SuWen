package test

import (
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/qqwry"
    "github.com/yhy0/SuWen/pkg/web"
    "github.com/yhy0/logging"
    "testing"
)

/**
   @author yhy
   @since 2024/5/21
   @desc //TODO
**/

func TestWeb(t *testing.T) {
    logging.Logger = logging.New(true, "", "SuWen", true)
    conf.Init()
    db.Init()
    qqwry.Init()
    // go TI.RunGithubAdvisories()
    web.Run("9088")
    
}
