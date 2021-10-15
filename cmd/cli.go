package main

import (
        "flag"
        "log"
        "fmt"
        "os"

        "github.com/HoMuChen/linenotify"
)

func main() {
        action := flag.String("action", "", "")
        token := flag.String("token", "", "")
        message := flag.String("message", "", "")

        flag.Parse()

        if *action == "" || (*action != "login" && *action != "send") {
                usage()
                return
        }

        client, err := newClient()
        if err != nil {
                log.Fatal(err)
        }

        if *action == "login" {
                token, err := client.Login()
                if err != nil {
                        log.Fatal(err)
                }

                fmt.Printf("access token: %s", token)
        }
        if *action == "send" {
                if err := client.Send(*token, *message); err != nil {
                        log.Fatal(err)
                }
        }
}

func usage() {
        fmt.Printf("usage: \n")
        fmt.Printf("\tcli -action %-6s\n", "login")
        fmt.Printf("\tcli -action %-6s -token {token} -message {message}\n", "send")
}

func newClient() (*linenotify.Client, error) {
        vars, err := loadEnvVar()
        if err != nil {
                return nil, err
        }

        client, err := linenotify.New(vars["LINE_CLIENT_ID"], vars["LINE_CLIENT_SECRET"], vars["LINE_CALLBACK"])
        if err != nil {
                return nil, err
        }

        return client, nil
}

func loadEnvVar() (map[string]string, error) {
        keys := []string{"LINE_CLIENT_ID", "LINE_CLIENT_SECRET", "LINE_CALLBACK"}
        ret := make(map[string]string)

        for _, key := range keys {
                val := os.Getenv(key)
                if val == "" {
                        return ret, fmt.Errorf("environment var %s should not be empty", key)
                }
                ret[key] = val
        }

        return ret, nil
}
