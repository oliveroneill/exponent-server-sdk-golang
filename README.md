# expo-server-sdk-go
Send push notifications to Expo apps using Go

## Installation
```
go get github.com/oliveroneill/exponent-server-sdk-golang/sdk
```

## Usage
```go
package main

import (
    "fmt"
    expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
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
            To: []expo.ExponentPushToken{pushToken},
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
    }
    
    // Validate responses
    if response.ValidateResponse() != nil {
        fmt.Println(response.PushMessage.To, "failed")
    }
}
```

## License
MIT
