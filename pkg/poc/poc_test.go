package poc

import (
    "context"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/notice"
    "github.com/yhy0/logging"
    "testing"
)

/**
  @author: yhy
  @since: 2024/5/27
  @desc: https://github.com/advisories
**/

func Test_Nuclei(t *testing.T) {
    logging.Logger = logging.New(true, "", "SuWen", true)
    conf.Init()
    db.Init()
    notice.InitPusher()
    FindNucleiPR(context.Background())
    
}

func Test_Afrog(t *testing.T) {
    logging.Logger = logging.New(true, "", "SuWen", true)
    conf.Init()
    db.Init()
    FindAfrog(context.Background())
    
}
