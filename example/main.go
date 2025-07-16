package main

import (
	"log"

	"github.com/cca2878/bsdkv3-go/sdk"
	"github.com/cca2878/bsdkv3-go/sdk/config"
	sdklog "github.com/cca2878/bsdkv3-go/sdk/log"
)

func main() {
	// Set log level for debugging
	sdklog.SetLevel(sdklog.LevelInfo)

	// Example 1: Basic login
	basicLoginExample()

	// Example 2: Login with custom configuration
	customConfigExample()

	// Example 3: Validator testing
	validatorExample()
}

func basicLoginExample() {
	log.Println("=== Basic Login Example ===")

	user := sdk.UserInfo{
		Username: "your_username",
		Password: "your_password",
	}

	client, err := sdk.NewBSdkV3Client(config.AppkeyPcr)
	if err != nil {
		log.Printf("创建客户端失败: %v\n", err)
		return
	}
	defer client.Close()

	accessKey, err := client.Login(user)
	if err != nil {
		log.Printf("登录失败: %v\n", err)
		return
	}

	log.Printf("登录成功！Access Key: %s\n", *accessKey)
}

func customConfigExample() {
	log.Println("\n=== Custom Configuration Example ===")

	// Create custom configuration
	cfg := config.NewConfig(
		config.WithAppId("custom_app_id"),
		config.WithPlatform("1"), // Different platform
	)

	client, err := sdk.NewBSdkV3Client(config.AppkeyPcr, sdk.WithConfig(cfg))
	if err != nil {
		log.Printf("创建客户端失败: %v\n", err)
		return
	}
	defer client.Close()

	log.Printf("客户端配置: AppId=%s, Platform=%s\n",
		client.GetConfig().RequestConfig.AppId,
		client.GetConfig().RequestConfig.Platform)
}

func validatorExample() {
	log.Println("\n=== Validator Example ===")

	validator := sdk.NewRemoteValidator()
	result, err := validator.Validate()
	if err != nil {
		log.Printf("验证失败: %v\n", err)
		return
	}

	log.Printf("验证成功!\n")
	log.Printf("Challenge: %s\n", result.Challenge)
	log.Printf("Gt: %s\n", result.Gt)
	log.Printf("GtUserId: %s\n", result.GtUserId)
	log.Printf("Validate: %s\n", result.Validate)
}
