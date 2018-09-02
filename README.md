# exponent-server-sdk-golang
[![Build Status](https://travis-ci.org/oliveroneill/exponent-server-sdk-golang.svg?branch=master)](https://travis-ci.org/oliveroneill/exponent-server-sdk-golang)

Exponent push notification go library based on [Expo's python sdk](https://github.com/exponent/exponent-server-sdk-python)

## Installation

```
go get github.com/oliveroneill/exponent-server-sdk-golang/sdk
```

## Usage

Use to send push notifications to Exponent Experiences from a Golang server.

[Full documentation](https://docs.expo.io/versions/latest/guides/push-notifications.html#http2-api) on the API is available if you want to dive into the details.

Example usage
```go
package main

import (
    "fmt"
    "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

func main() {
    // To check the token is valid
    pushToken, err := expo.NewExponentPushToken("ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]")
    if err != nil {
        panic(err)
    }

    // Create a new Expo SDK client
    client := expo.NewPushClient(nil)

    // Publish message
    response, err := client.Publish(
        &expo.PushMessage{
            To: pushToken,
            Body: "This is a test notification",
            Data: map[string]string{"withSome": "data"},
            Sound: "default",
            Title: "Notification Title",
            Priority: expo.DefaultPriority,
        },
    )
    // Check errors
    if err != nil {
        panic(err)
        return
    }
    // Validate responses
    if response.ValidateResponse() != nil {
        fmt.Println(response.PushMessage.To, "failed")
    }
}
```

## TODO

  * Need to add more unit tests

## See Also

  * https://github.com/exponent/exponent-server-sdk-ruby
  * https://github.com/exponent/exponent-server-sdk-python
  * https://github.com/exponent/exponent-server-sdk-node
