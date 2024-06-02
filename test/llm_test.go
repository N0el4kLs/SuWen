package test

import (
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/SuWen/pkg/llm"
    "github.com/yhy0/logging"
    "testing"
)

/**
  @author: yhy
  @since: 2024/5/22
  @desc: //TODO
**/

func TestOpenAI(t *testing.T) {
    logging.Logger = logging.New(true, "", "llm", true)
    conf.Init()
    db.Init()
    llm.ChatGPT("")
}
