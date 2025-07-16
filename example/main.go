package main

import (
	"bsdkv3-go/sdk"
	"bsdkv3-go/sdk/config"
	"bsdkv3-go/sdk/log"
)

func testValidator() {
	log.SetLevel(log.LevelDebug)
	ret, err := bsdkv3.NewRemoteValidator().Validate()
	if err != nil {
		log.Error("验证失败！")
		log.Error("错误信息：%s", err.Error())
		return
	}
	// 直接使用结构化的结果
	log.Info("验证成功！")
	log.Info("Challenge: %s", ret.Challenge)
	log.Info("Gt: %s", ret.Gt)
	log.Info("GtUserId: %s", ret.GtUserId)
	log.Info("Validate: %s", ret.Validate)
}

func testLogin() {
	log.SetLevel(log.LevelDebug)
	user := bsdkv3.UserInfo{
		Username: "your_username",
		Password: "your_password",
	}
	client, err := bsdkv3.NewClient(config.AppkeyPcr)
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
	log.Info("登录成功！")
	log.Info("登录结果: %v", ret)
	client.Close()
}

func main() {
	testLogin()

}
