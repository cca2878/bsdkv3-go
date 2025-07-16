# BSdk V3 Client Go

[![Go Version](https://img.shields.io/badge/go-1.24%2B-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-CC%20BY--NC--SA%204.0-green.svg)](https://creativecommons.org/licenses/by-nc-sa/4.0/)

A Go SDK client implementation for BSdk V3 authentication and validation services.

## Features

- 🔐 User authentication with password encryption
- 🤖 Captcha handling with remote validation
- ⚙️ Configurable client options
- 📝 Comprehensive logging
- 🔄 Automatic host configuration updates
- 🚀 Simple and easy-to-use API

## Installation

```bash
go get github.com/cca2878/bsdkv3-go
```

## Quick Start

```go
package main

import (
    "log"
    
    "bsdkv3-go/sdk"
    "bsdkv3-go/sdk/config"
)

func main() {
    // Create a new client
    client, err := bsdkv3.NewClient(config.AppkeyPcr)
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }
    defer client.Close()

    // Prepare user info
    user := bsdkv3.UserInfo{
        Username: "your_username",
        Password: "your_password",
        Platform: "android",
        Channel:  "bilibili",
    }

    // Login
    account, err := client.Login(user)
    if err != nil {
        log.Fatal("Login failed:", err)
    }

    log.Printf("Login successful! UID: %s, AccessKey: %s", account.Uid, account.AccessKey)
}
```

## Configuration

The SDK supports various configuration options:

```go
import "bsdkv3-go/sdk/config"

// Create custom configuration
cfg := config.NewConfig(
    config.WithRequestTimeout(10 * time.Second),
    config.WithAppId("your_app_id"),
    config.WithGameId("your_game_id"),
    // ... more options
)

// Create client with custom config
client, err := bsdkv3.NewClient(appKey, bsdkv3.WithConfig(cfg))
```

## Available App Keys

- `config.AppkeyPcr` - Princess Connect Re:Dive

## API Reference

### Client

#### `NewClient(appKey string, options ...ClientOption) (*Client, error)`

Creates a new SDK client instance.

#### `Login(user UserInfo) (*SdkAccount, error)`

Authenticates a user and returns account information.

#### `Close()`

Closes the client and cancels all pending requests.

### Types

#### `UserInfo`

```go
type UserInfo struct {
    Username string // User's username
    Password string // User's password
    Platform string // Platform identifier
    Channel  string // Channel identifier
}
```

#### `SdkAccount`

```go
type SdkAccount struct {
    Uid       string // User ID
    AccessKey string // Access token
    Platform  string // Platform identifier
    Channel   string // Channel identifier
}
```

## Logging

The SDK includes a built-in logging system with configurable levels:

```go
import "bsdkv3-go/sdk/log"

// Set log level
log.SetLevel(log.LevelDebug) // Debug, Info, Warn, Error

// Set custom output
file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
log.SetOutput(file)
```

## Error Handling

The SDK provides comprehensive error handling with detailed error messages. All errors include context about what operation failed and why.

## Contributing

This project is for personal use and learning purposes. 

## License

This project is licensed under the CC BY-NC-SA 4.0 License - see the [LICENSE](https://creativecommons.org/licenses/by-nc-sa/4.0/) for details.

## Disclaimer

This SDK is for educational and personal use only. Users are responsible for complying with the terms of service of the target platform.
