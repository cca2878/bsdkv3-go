# BSdk V3 Client Go

[![Go Reference](https://pkg.go.dev/badge/github.com/cca2878/bsdkv3-go.svg)](https://pkg.go.dev/github.com/cca2878/bsdkv3-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/cca2878/bsdkv3-go)](https://goreportcard.com/report/github.com/cca2878/bsdkv3-go)

## Overview

A Go SDK client for authentication and validation services. This package provides a clean, idiomatic Go interface for authentication workflows including captcha handling.

## Features

- **Easy to use**: Simple API for authentication
- **Captcha support**: Automatic captcha handling via remote validation
- **Configurable**: Flexible configuration options
- **Context support**: Proper context handling for request cancellation
- **Type safe**: Full type safety with structured request/response models

## Installation

```bash
go get github.com/cca2878/bsdkv3-go
```

## Quick Start

```go
package main

import (
    "log"
    
    "github.com/cca2878/bsdkv3-go/sdk"
    "github.com/cca2878/bsdkv3-go/sdk/config"
)

func main() {
    // Create a new client
    client, err := sdk.NewBSdkV3Client(config.AppkeyPcr)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Login with user credentials
    user := sdk.UserInfo{
        Username: "your_username",
        Password: "your_password",
    }
    
    accessKey, err := client.Login(user)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Login successful, access key: %s", *accessKey)
}
```

## Configuration

The SDK supports flexible configuration through the config package:

```go
import "github.com/cca2878/bsdkv3-go/sdk/config"

// Use default configuration
cfg := config.NewDefaultConfig()

// Or customize specific options
cfg := config.NewConfig(
    config.WithRequestTimeout(10 * time.Second),
    config.WithAppId("your_app_id"),
)

// Use custom config with client
client, err := sdk.NewBSdkV3Client(appKey, sdk.WithConfig(cfg))
```

## API Documentation

For detailed API documentation, visit [pkg.go.dev](https://pkg.go.dev/github.com/cca2878/bsdkv3-go).

## Example

See the [example directory](./example) for a complete working example.

## License

CC BY-NC-SA 4.0

## Contributing

This project is primarily for personal use and learning purposes. Issues and pull requests are welcome.
