package conf

import (
    "bytes"
    "github.com/google/go-github/v62/github"
    "github.com/spf13/viper"
    "github.com/yhy0/SuWen/pkg/util"
    "github.com/yhy0/logging"
    "net/http"
    "path"
    "time"
)

type Config struct {
    DbConfig           DbConfig           `json:"dbConfig"`
    WatchVulnAppConfig WatchVulnAppConfig `json:"watchVulnAppConfig"`
    LLMConfig          LLMConfig          `json:"llmConfig"`
    GithubToken        string             `json:"githubToken"`
}

type DbConfig struct {
    Host     string `json:"host"`
    Password string `json:"password"`
    Port     string `json:"port"`
    User     string `json:"user"`
    Database string `json:"database"`
    Timeout  string `json:"timeout"`
}

type WatchVulnAppConfig struct {
    Sources         string `json:"sources"`           // 启用哪些源
    Interval        string `json:"interval"`          // 每次获取情报的时间间隔
    EnableCVEFilter bool   `json:"enable_cve_filter"` // 启用过滤器，该过滤器来自具有相同cve id的多个来源的vulns将仅发送一次
    NoGithubSearch  bool   `json:"no_github_search"`  // 不要搜索github repo并为每个cve vuln拉取请求
    NoStartMessage  bool   `json:"no_start_message"`  // 服务器启动时禁用问候消息
    NoFilter        bool   `json:"no_filter"`         // 忽略有价值的过滤器并推送所有发现的情报
    DiffMode        bool   `json:"diff_mode"`         // 跳过初始化vuln DB，推送新vuln然后退出
}

type LLMConfig struct {
    Model string `json:"model"`
    Token string `json:"token"`
}

// Init 加载配置
func Init() {
    // 配置文件路径 当前文件夹 + SScan_config.yaml
    configFile := path.Join("./", ConfigFileName)
    
    // 检测配置文件是否存在
    if !util.Exists(configFile) {
        WriteYamlConfig(configFile)
        logging.Logger.Infof("%s not find, Generate profile.", configFile)
    } else {
        logging.Logger.Infoln("Load profile ", configFile)
    }
    ReadYamlConfig(configFile)
    
    tr := http.DefaultTransport.(*http.Transport).Clone()
    tr.Proxy = http.ProxyFromEnvironment
    GithubClient = github.NewClient(&http.Client{
        Timeout:   time.Second * 10,
        Transport: tr,
    }).WithAuthToken(GlobalConfig.GithubToken)
    
}

func ReadYamlConfig(configFile string) {
    // 加载config
    viper.SetConfigType("yaml")
    viper.SetConfigFile(configFile)
    
    err := viper.ReadInConfig()
    if err != nil {
        logging.Logger.Fatalf("setting.Setup, fail to read '%s': %+v", ConfigFileName, err)
    }
    err = viper.Unmarshal(&GlobalConfig)
    
    if err != nil {
        logging.Logger.Fatalf("setting.Setup, fail to parse '%s', check format: %v", ConfigFileName, err)
    }
}

func WriteYamlConfig(configFile string) {
    // 生成默认config
    viper.SetConfigType("yaml")
    err := viper.ReadConfig(bytes.NewBuffer(defaultYamlByte))
    if err != nil {
        logging.Logger.Fatalf("setting.Setup, fail to read default config bytes: %v", err)
    }
    // 写文件
    err = viper.SafeWriteConfigAs(configFile)
    if err != nil {
        logging.Logger.Fatalf("setting.Setup, fail to write '%s': %v", ConfigFileName, err)
    }
}
