# zapmixin

[![Go Report Card](https://goreportcard.com/badge/github.com/xwjdsh/zapmixin)](https://goreportcard.com/report/github.com/xwjdsh/zapmixin)

Hook for sending events zap logger to mixin.

## Usage

```go
package main

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/xwjdsh/zapmixin"
	"go.uber.org/zap"
)

func main() {
	client, err := mixin.NewFromKeystore(&mixin.Keystore{
		ClientID:   "<client_id>",
		SessionID:  "<session_id>",
		PrivateKey: "<private_key>",
		PinToken:   "<pin_token>",
		Scope:      "<scope>",
	})
	if err != nil {
		panic(err)
	}

	conversations := []string{"<conversation_id>"}
	h, err := zapmixin.New(client, conversations)

	logger, _ := zap.NewProduction()
	logger.WithOptions(zap.Hooks(h.Hook()))

	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
}

```

## Parameters

#### Required

- mixin.Client (by https://github.com/fox-one/mixin-sdk-go)
- Conversation IDs

#### Optional

- WithThresholdLevel - Level threshold, the default level is WARN.
- WithFixedLevel - Only the given level will send the message.
- WithSync - The default is to send messages to mixin asynchronously.
- WithFormatter - Custom formatter.
- WithFilter - Filters are applied to messages to determine if any entry should not be send out.
- WithAfter - When the message is sent, it will be called.

## Installation

```
go get github.com/xwjdsh/zapmixin
```

## Credits

Based on [zaptelegram](https://github.com/strpc/zaptelegram).
