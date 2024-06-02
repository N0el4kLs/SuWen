package db

import "gorm.io/gorm"

/**
  @author: yhy
  @since: 2024/5/21
  @desc: //TODO
**/

// Pusher 通知推送
type Pusher struct {
    gorm.Model
    Id      int    `gorm:"primary_key" json:"id"`
    Source  string `json:"source"`
    Webhook string `json:"webhook"`
    Token   string `json:"token"`
    Secret  string `json:"secret"`
}

func GetPusher() (pushers []Pusher) {
    GlobalDB.Model(&Pusher{}).Order("id asc").Find(&pushers)
    
    return
}
