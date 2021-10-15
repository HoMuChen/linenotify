package linenotify

import (
        "testing"
        "net/http"
        "strings"
)

func TestSend(t *testing.T) {
        c, _ := New(ID, SECRET, CB)

        message := "this is a message for test"
        req := c.makeNotifyRequest("token", message)

        if req.Method != http.MethodPost {
                t.Error()
        }
        if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
                t.Error()
        }
        if !strings.Contains(req.Header.Get("Authorization"), "token") {
                t.Error()
        }

        err := req.ParseForm()
        if err != nil {
                t.Error(err)
        }
        if req.Form.Get("message") != message {
                t.Errorf("form data with key message should be %s, but got: %s", message, req.Form.Get("message"))
        }
}
