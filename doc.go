/*
Package bsdkv3-go provides a Go client SDK for authentication and validation services.

This SDK offers a clean, idiomatic Go interface for handling authentication
workflows, including automatic captcha validation and secure credential handling.

# Installation

	go get github.com/cca2878/bsdkv3-go

# Quick Start

	package main

	import (
		"log"

		"github.com/cca2878/bsdkv3-go/sdk"
		"github.com/cca2878/bsdkv3-go/sdk/config"
	)

	func main() {
		client, err := sdk.NewBSdkV3Client(config.AppkeyPcr)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		user := sdk.UserInfo{
			Username: "your_username",
			Password: "your_password",
		}

		accessKey, err := client.Login(user)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Login successful: %s", *accessKey)
	}

# Features

  - Easy-to-use client interface
  - Automatic captcha handling via remote validation
  - Flexible configuration system
  - Proper context handling for request cancellation
  - Type-safe request/response models
  - Comprehensive error handling with custom error types

# Error Handling

The SDK provides specific error types for different failure scenarios:

  - *sdk.AuthError: Authentication failures with error codes
  - *sdk.ConfigError: Configuration-related errors
  - *sdk.CaptchaError: Captcha validation failures

Example error handling:

	accessKey, err := client.Login(user)
	if err != nil {
		var authErr *sdk.AuthError
		if errors.As(err, &authErr) {
			log.Printf("Authentication failed [%s]: %s", authErr.Code, authErr.Message)
		}
		return err
	}

# Configuration

The client supports extensive customization through the config package:

	cfg := config.NewConfig(
		config.WithRequestTimeout(10 * time.Second),
		config.WithAppId("custom_app_id"),
		config.WithPlatform("1"),
	)
	client, err := sdk.NewBSdkV3Client(appKey, sdk.WithConfig(cfg))

For complete documentation, visit https://pkg.go.dev/github.com/cca2878/bsdkv3-go
*/
package main