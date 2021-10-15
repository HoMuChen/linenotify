package linenotify

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "net/url"
        "strings"
)


func (c *client) Send(token, message string) error {
        req := c.makeNotifyRequest(token, message)
        res, err := http.DefaultClient.Do(req)
        if err != nil {
                return fmt.Errorf("fail to post to notify api with err: %v", err)
        }
        defer res.Body.Close()

        data, err := ioutil.ReadAll(res.Body)
        if err != nil {
                return fmt.Errorf("fail to read response body  with err %v", err)
        }

        var buf struct {
                Status          int     `json:status`
                Message         string  `json:message`
        }
        if err := json.Unmarshal(data, &buf); err != nil {
                return fmt.Errorf("fail to json parse body with err %v", err)
        }
        if buf.Status >= 400 {
                return fmt.Errorf("fail to send message with status: %v and res message: %s", buf.Status, buf.Message)
        }

        return nil
}

func (c *client) makeNotifyRequest(token, message string) *http.Request {
        values := url.Values{}
        values.Add("message", message)
        body := values.Encode()

        req, _ := http.NewRequest("POST", NOTIFY_HOST, strings.NewReader(body))
        req.Header.Add("Content-type", "application/x-www-form-urlencoded")
        req.Header.Add("Authorization", "Bearer " + token)

        return req
}
