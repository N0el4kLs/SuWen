package notice

import (
    "encoding/json"
    "fmt"
    "github.com/go-resty/resty/v2"
    "github.com/pkg/errors"
    "github.com/yhy0/logging"
)

type PlusMessage struct {
    Token    string `json:"token"`
    Title    string `json:"title" describe:"消息标题"`
    Content  string `json:"content" describe:"具体消息内容，根据不同template支持不同格式"`
    Template string `json:"template" describe:"发送消息模板"`
}

type PlusResponse struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
    Data string `json:"data"`
}

var _ = TextPusher(&Plus{})

type Plus struct {
    token string
}

func NewPushPlus(token string) TextPusher {
    return &Plus{
        token: token,
    }
}

func (r *Plus) Send(message PlusMessage) (response *PlusResponse, error error) {
    res := &PlusResponse{}
    message.Token = r.token
    
    if len(message.Token) == 0 {
        return res, errors.New("invalid token")
    }
    
    result, err := resty.New().R().SetBody(message).SetHeader("Content-Type", "application/json").Post("https://www.pushplus.plus/send")
    
    if err != nil {
        return res, errors.New(fmt.Sprintf("请求失败：%s", err.Error()))
    }
    err = json.Unmarshal(result.Body(), res)
    if err != nil {
        return res, errors.New("json 格式化数据失败")
    }
    if res.Code != 200 {
        return res, errors.New(res.Msg)
    }
    return res, nil
}

func (d *Plus) PushText(s string) error {
    logging.Logger.Infof("sending text %s", s)
    message := PlusMessage{
        Title:    "",
        Content:  s,
        Template: "txt",
    }
    
    _, err := d.Send(message)
    if err != nil {
        return errors.Wrap(err, "push-plus")
    }
    return nil
}

func (d *Plus) PushMarkdown(title, content string) error {
    logging.Logger.Infof("sending markdown %s", title)
    message := PlusMessage{
        Title:    title,
        Content:  content,
        Template: "markdown",
    }
    
    _, err := d.Send(message)
    if err != nil {
        return errors.Wrap(err, "push-plus")
    }
    return nil
}
