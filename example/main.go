package main

import (
	"bsdkv3/sdk"
	"bsdkv3/sdk/config"
	"bsdkv3/sdk/log"
	"fmt"
)

func testValidator() {
	log.SetLevel(log.LevelDebug)
	ret, err := sdk.NewRemoteValidator().Validate()
	if err != nil {
		log.Error("验证失败！")
		log.Error("错误信息：%s", err.Error())
		return
	}
	// 直接使用结构化的结果
	fmt.Printf("验证成功！\n")
	fmt.Printf("Challenge: %s\n", ret.Challenge)
	fmt.Printf("Gt: %s\n", ret.Gt)
	fmt.Printf("GtUserId: %s\n", ret.GtUserId)
	fmt.Printf("Validate: %s\n", ret.Validate)
}

func testLogin() {
	log.SetLevel(log.LevelDebug)
	user := sdk.UserInfo{
		Username: "your_username",
		Password: "your_password",
	}
	client, err := sdk.NewBSdkV3Client(config.AppkeyPcr)
	if err != nil {
		log.Error("创建客户端失败！")
		log.Error("错误信息：%s", err.Error())
		return
	}
	ret, err := client.Login(user)
	if err != nil {
		log.Error("登录失败！")
		log.Error("错误信息：%s", err.Error())
		return
	}
	log.Info("登录成功！" + *ret)
	client.Close()
}

func main() {
	testLogin()
}
