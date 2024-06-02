package test

import (
    "github.com/yhy0/SuWen/pkg/TI"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/web"
    "github.com/yhy0/logging"
    "testing"
)

/**
  @author: yhy
  @since: 2024/5/22
  @desc: //TODO
**/

func TestApp(t *testing.T) {
    logging.Logger = logging.New(true, "", "SuWen", true)
    conf.Init()
    db.Init()
    
    go func() {
        err := TI.RunPR()
        if err != nil {
            logging.Logger.Errorln(err)
        }
    }()
    web.Run("9088")
}
