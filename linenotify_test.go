package linenotify

import (
        "testing"
)

const (
        ID      = "id"
        SECRET  = "secret"
        CB      = "http://localhost:5000/callback"
)

func TestNewWithoutSecret(t *testing.T) {
        _, err := New(ID, "", "")

        if err == nil {
                t.Errorf("error should not nil")
        }
}

func TestNew(t *testing.T) {
        _, err := New(ID, SECRET, CB)

        if err != nil {
                t.Errorf("error should be nil, but got: %v", err)
        }
}

func TestAuthUrl(t *testing.T) {
        c, _ := New(ID, SECRET, CB)

        url := c.AuthUrl()
        expected := "https://notify-bot.line.me/oauth/authorize?client_id=id&redirect_uri=http%3A%2F%2Flocalhost%3A5000%2Fcallback&response_type=code&scope=notify&state=123"

        if url != expected {
                t.Errorf("expected: %s, but got: %s", expected, url)
        }
}

func TestMakeTokenBody(t *testing.T) {
        c, _ := New(ID, SECRET, CB)

        body := c.makeTokenBody("abctoken")
        expected := "client_id=id&client_secret=secret&code=abctoken&grant_type=authorization_code&redirect_uri=http%3A%2F%2Flocalhost%3A5000%2Fcallback"

        if body != expected {
                t.Errorf("expected: %s, but got: %s", expected, body)
        }
}
