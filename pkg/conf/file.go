package conf

var defaultYamlByte = []byte(`
# 数据库配置
dbConfig:
  host: "127.0.0.1"
  password: ""
  port: "3306"
  user: "root"
  database: "suwen"
  # 数据库连接超时时间
  timeout: "3s"

# watchVuln 配置
watchVulnAppConfig:
  sources: "avd,nox,oscs,threatbook,seebug,struts2,kev"
  interval: 30m
  enable_cve_filter: true
  no_github_search: false
  no_start_message: false
  no_filter: false
  diff_mode: false
llmConfig:
  # gpt4 免费的每天只有 3 条, gpt3 每天 100 条
  model: gpt-3.5-turbo
  token: ""
githubToken: ""
`)
