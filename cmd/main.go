package cmd

import (
    "github.com/urfave/cli/v2"
    "github.com/yhy0/SuWen/pkg/TI"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/notice"
    "github.com/yhy0/SuWen/pkg/poc"
    "github.com/yhy0/SuWen/pkg/qqwry"
    "github.com/yhy0/SuWen/pkg/web"
    "github.com/yhy0/logging"
    "os"
)

/**
  @author: yhy
  @since: 2024/5/21
  @desc: //TODO
**/

func RunApp() {
    logging.Logger = logging.New(true, "", "SuWen", true)
    app := cli.NewApp()
    app.Name = "素问"
    app.Usage = conf.Website
    app.Version = conf.Version
    
    app.Flags = []cli.Flag{
        &cli.StringFlag{
            // 后端端口
            Name:    "port",
            Aliases: []string{"p"},
            Value:   conf.DefaultWebPort,
            Usage:   "web server `PORT`",
        },
    }
    
    app.Action = RunServer
    
    err := app.Run(os.Args)
    
    if err != nil {
        logging.Logger.Fatalf("cli.RunApp err: %+v", err)
        return
    }
}

func RunServer(ctx *cli.Context) error {
    // config 必须最先加载
    conf.Init()
    db.Init()
    notice.InitPusher()
    qqwry.Init()
    
    go func() {
        go TI.RunGithubAdvisories()
        go poc.Run()
        
        err := TI.RunPR()
        if err != nil {
            logging.Logger.Errorln(err)
        }
        
    }()
    web.Run(ctx.String("port"))
    return nil
}
