# Linenotify
A command tool for line-notify service, including login to get a token and send message.

## Usage
```go
import (
        "log"

        "github.com/HoMuChen/linenotify"
)

func main() {
        client, err := linenotify("id", "secret", "callback")
        if err != nil {
                log.Fatal(err)
        }

        toten, err := client.Login()
        if err != nil {
                log.Fatal(err)
        }

        if err := client.send("token", "hello"); err != nil {
                log.Fatal(err)
        }
}
```

## Command line

### Build
```sh
go build ./cmd/cli.go
```

### Configuration
```sh
export LINE_CLIENT_ID="your client id"
export LINE_CLIENT_SECRET="your client secret"
export LINE_CALLBACK="callback url"
```

### Run
* Login
  ```sh
  ./cli -action login
  ```

* Send message
  ```sh
  ./cli -action send -token {token} -messgae {message}
  ```
