package notice

import (
    "github.com/yhy0/SuWen/pkg/db"
    "github.com/yhy0/logging"
    "strings"
)

/**
   @author yhy
   @since 2024/5/30
   @desc //TODO
**/

var Text TextPusher
var Raw RawPusher

func InitPusher() {
    var textPusher []TextPusher
    var rawPusher []RawPusher
    
    pushers := db.GetPusher()
    for _, pusher := range pushers {
        if pusher.Source == "ding" && pusher.Token != "" && pusher.Secret != "" {
            textPusher = append(textPusher, NewDingDing(pusher.Token, pusher.Secret))
        } else if pusher.Source == "lark" && pusher.Token != "" && pusher.Secret != "" {
            textPusher = append(textPusher, NewLark(pusher.Token, pusher.Secret))
        } else if pusher.Source == "wx" && pusher.Token != "" {
            textPusher = append(textPusher, NewWechatWork(pusher.Token))
        } else if pusher.Source == "webhook" && pusher.Webhook != "" {
            rawPusher = append(rawPusher, NewWebhook(pusher.Webhook))
        } else if pusher.Source == "lanxin" && pusher.Webhook == "lanxin" && pusher.Token != "" && pusher.Secret != "" {
            textPusher = append(textPusher, NewLanxin(pusher.Webhook, pusher.Token, pusher.Secret))
        } else if pusher.Source == "bark" && pusher.Webhook != "" {
            deviceKeys := strings.Split(pusher.Webhook, "/")
            deviceKey := deviceKeys[len(deviceKeys)-1]
            url := strings.Replace(pusher.Webhook, deviceKey, "push", -1)
            textPusher = append(textPusher, NewBark(url, deviceKey))
        } else if pusher.Source == "serverChan" && pusher.Token != "" {
            textPusher = append(textPusher, NewServerChan(pusher.Token))
        } else if pusher.Source == "pushPlus" && pusher.Token != "" {
            textPusher = append(textPusher, NewPushPlus(pusher.Token))
        } else if pusher.Source == "telegram" && pusher.Token != "" && pusher.Secret != "" {
            tgPusher, err := NewTelegram(pusher.Token, pusher.Secret)
            if err != nil {
                logging.Logger.Errorln(err)
            }
            textPusher = append(textPusher, tgPusher)
        }
    }
    Text = MultiTextPusher(textPusher...)
    Raw = MultiRawPusher(rawPusher...)
}
