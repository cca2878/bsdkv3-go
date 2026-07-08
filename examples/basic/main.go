package main

import (
	"context"
	"fmt"
	"time"

	bsdkv3 "github.com/cca2878/bsdkv3-go"
)

func main() {
	fmt.Println("============== bsdkv3 login example ==============")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// NewClient 会隐式完成 bootstrap（拉取登录 host 列表 + RSA 登录密钥），fail-fast。
	client, err := bsdkv3.NewClient(ctx, bsdkv3.AppkeyPcr)
	if err != nil {
		fmt.Println("创建客户端失败:", err)
		return
	}

	// 执行登录并自动处理验证码（geetest 求解走独立 HTTP 客户端）。
	acc, err := client.Auth.Login(ctx, bsdkv3.UserInfo{
		Username: "example_user",
		Password: "example_password",
	})
	if err != nil {
		fmt.Println("登录失败:", err)
		return
	}

	fmt.Printf("登录成功 Uid=%s AccessKey=%s\n", acc.Uid, acc.AccessKey)
}
