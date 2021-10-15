package linenotify

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "net/url"
        "os/exec"
        "strings"
        "runtime"
)

const (
        LINE_AUTH_PATH          = "https://notify-bot.line.me/oauth/authorize"
        LINE_TOKEN_PATH         = "https://notify-bot.line.me/oauth/token"
        NOTIFY_HOST             = "https://notify-api.line.me/api/notify"
)

type client struct {
        id              string
        secret          string
        callback        *url.URL
        token           string
        done            chan error
}

func New(id, secret, callback string) (*client, error) {
        if id == "" || secret == "" {
                return nil, fmt.Errorf("client id: %s and secret: %s can not be empty", id, secret)
        }

        cb, err := url.Parse(callback)
        if err != nil {
                return nil, fmt.Errorf("fail to pars callback url with err: %v", err)
        }

        c := &client{
                id:             id,
                secret:         secret,
                callback:       cb,
                token:          "",
                done:           make(chan error),
        }

        return c, nil
}

func (c *client) Login() (string, error) {
        if err := openbrowser(c.AuthUrl()); err != nil {
                return "", fmt.Errorf("fail to open browser, %v", err)
        }

        if err := c.waitForCallback(); err != nil {
                return "", fmt.Errorf("fail to start server waiting for callback code to exchange access toekn: %v", err)
        }

        return c.token, nil
}

func (c *client) AuthUrl() string {
        values := url.Values{}
        values.Add("client_id", c.id)
        values.Add("response_type", "code")
        values.Add("scope", "notify")
        values.Add("redirect_uri", c.callback.String())
        values.Add("state", "123")
        query := values.Encode()

        url := LINE_AUTH_PATH + "?" + query

        return url
}

func (c *client) makeTokenBody(code string) string {
        values := url.Values{}
        values.Add("client_id", c.id)
        values.Add("client_secret", c.secret)
        values.Add("redirect_uri", c.callback.String())
        values.Add("code", code)
        values.Add("grant_type", "authorization_code")
        body := values.Encode()

        return body
}

func (c *client) GetToken(code string) (string, error) {
        body := c.makeTokenBody(code)
        req, _ := http.NewRequest("POST", LINE_TOKEN_PATH, strings.NewReader(body))
        req.Header.Add("Content-type", "application/x-www-form-urlencoded")

        res, err := http.DefaultClient.Do(req)
        if err != nil {
                return "", fmt.Errorf("fail to get token with err: %v", err)
        }
        defer res.Body.Close()

        data, err := ioutil.ReadAll(res.Body)
        if err != nil {
                return "", fmt.Errorf("fail to read response body  with err %v", err)
        }

        var buf struct {
                Access_token    string  `json:access_token`
                Status          int     `json:status`
                Message         string  `json:message`
        }
        if err := json.Unmarshal(data, &buf); err != nil {
                return "", fmt.Errorf("fail to json parse body with err %v", err)
        }

        if buf.Status >= 400 {
                return "", fmt.Errorf("fail to get token with status: %v and message: %s", buf.Status, buf.Message)
        }

        return buf.Access_token, nil
}

func (c *client) handler(w http.ResponseWriter, r *http.Request) {
        code := r.URL.Query().Get("code")
        token, err := c.GetToken(code)
        if err != nil {
                c.done <- err
                return
        }

        w.Write([]byte(token))

        c.token = token
        c.done <- nil
}

func (c *client) waitForCallback() error {
        server := &http.Server{
                Addr:           c.callback.Host,
                Handler:        http.HandlerFunc(c.handler),
        }

        go server.ListenAndServe()
        defer server.Close()

        return <-c.done
}

func openbrowser(url string) error {
        var err error

        switch runtime.GOOS {
        case "linux":
                err = exec.Command("xdg-open", url).Start()
        case "windows":
                err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
        case "darwin":
                err = exec.Command("open", url).Start()
        default:
                err = fmt.Errorf("unsupported platform")
        }

        return err
}
