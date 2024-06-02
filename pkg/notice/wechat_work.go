package notice

import (
    "github.com/pkg/errors"
    wxworkbot "github.com/vimsucks/wxwork-bot-go"
    "github.com/yhy0/logging"
)

var _ = TextPusher(&WechatWork{})

type WechatWork struct {
    client *wxworkbot.WxWorkBot
}

func NewWechatWork(botKey string) TextPusher {
    return &WechatWork{
        client: wxworkbot.New(botKey),
    }
}

func (d *WechatWork) PushText(s string) error {
    // fixme: wxworkbot 不支持 text 类型
    logging.Logger.Infof("sending text %s", s)
    msg := wxworkbot.Markdown{Content: s}
    err := d.client.Send(msg)
    if err != nil {
        return errors.Wrap(err, "wechat-work")
    }
    return nil
}

func (d *WechatWork) PushMarkdown(title, content string) error {
    logging.Logger.Infof("sending markdown %s", title)
    msg := wxworkbot.Markdown{Content: content}
    err := d.client.Send(msg)
    if err != nil {
        return errors.Wrap(err, "wechat-work")
    }
    return nil
}
