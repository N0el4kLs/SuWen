package notice

import (
    "github.com/pkg/errors"
    "github.com/rayepeng/serverchan"
    "github.com/yhy0/logging"
)

var _ = TextPusher(&ServerChan{})

type ServerChan struct {
    client *serverchan.ServerChan
}

func NewServerChan(botKey string) TextPusher {
    return &ServerChan{
        client: serverchan.NewServerChan(botKey),
    }
}

func (d *ServerChan) PushText(s string) error {
    logging.Logger.Infof("sending text %s", s)
    _, err := d.client.Send("", s)
    if err != nil {
        return errors.Wrap(err, "server-chan")
    }
    return nil
}

func (d *ServerChan) PushMarkdown(title, content string) error {
    logging.Logger.Infof("sending markdown %s", title)
    _, err := d.client.Send(title, content)
    if err != nil {
        return errors.Wrap(err, "server-chan")
    }
    return nil
}
