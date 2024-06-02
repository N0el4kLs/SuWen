package conf

import (
    "github.com/google/go-github/v62/github"
)

/**
   @author yhy
   @since 2024/6/2
   @desc //TODO
**/

var GlobalConfig *Config

var GithubClient *github.Client

const ConfigFileName = "SuWen.yaml"

const DefaultWebPort = "9088"

var LastCheckTime = make(map[string]string)
