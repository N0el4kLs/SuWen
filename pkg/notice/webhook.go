package notice

import (
    "bytes"
    "encoding/json"
    "github.com/pkg/errors"
    "github.com/yhy0/logging"
    "io"
    "net/http"
)

var _ = RawPusher(&Webhook{})

type Webhook struct {
    url    string
    client *http.Client
}

func NewWebhook(url string) RawPusher {
    return &Webhook{
        url:    url,
        client: &http.Client{},
    }
}

func (m *Webhook) PushRaw(r *RawMessage) error {
    logging.Logger.Infof("sending webhook data %s, %v", r.Type, r.Content)
    postBody, _ := json.Marshal(r)
    resp, err := m.doPostRequest(m.url, "application/json", postBody)
    if err != nil {
        return err
    }
    logging.Logger.Infof("raw response from server: %s", string(resp))
    return nil
}

func (m *Webhook) doPostRequest(url string, contentType string, body []byte) ([]byte, error) {
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
    if err != nil {
        return nil, errors.Wrap(err, "create request")
    }
    
    req.Header.Set("Content-Type", contentType)
    
    resp, err := m.client.Do(req)
    if err != nil {
        return nil, errors.Wrap(err, "send request")
    }
    defer resp.Body.Close()
    
    respBody, err := io.ReadAll(resp.Body) // 使用 io.ReadAll 替代 ioutil.ReadAll。
    if err != nil {
        return nil, errors.Wrap(err, "read body")
    }
    return respBody, nil
}
