package llm

import (
    "context"
    "fmt"
    "github.com/sashabaranov/go-openai"
    "github.com/yhy0/SuWen/pkg/conf"
    "github.com/yhy0/logging"
    "strings"
)

/**
   @author yhy
   @since 2024/5/29
   @desc 感谢 https://github.com/chatanywhere/GPT_API_free 提供的免费 api
**/

func ChatGPT(msg string) string {
    config := openai.DefaultConfig(conf.GlobalConfig.LLMConfig.Token)
    
    config.BaseURL = "https://api.chatanywhere.com.cn/v1"
    client := openai.NewClientWithConfig(config)
    
    resp, err := client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: openai.GPT3Dot5Turbo,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleUser,
                    Content: fmt.Sprintf("请帮我把以下内容翻译成中文，同时对输出的 markdown 语法和链接信息 前后空一格显示，并且只返回你处理后的内容,不要增加任何无关输出\n %s", msg),
                },
            },
        },
    )
    
    if err != nil {
        logging.Logger.Errorf("ChatCompletion error: %v", err)
        return ""
    }
    
    res := strings.ReplaceAll(resp.Choices[0].Message.Content, "```html", "")
    return strings.TrimSuffix(res, "```")
}
