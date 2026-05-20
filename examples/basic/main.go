package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cca2878/bsdkv3-go"
)

func testLoginR() {
	// 创建日志记录器，使用导出的 LogLevelDebug 级别
	logger := bsdkv3.NewStdLogger(os.Stdout, bsdkv3.LogLevelDebug)

	user := bsdkv3.UserInfo{
		Username: "example_user",
		Password: "example_password",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := bsdkv3.NewClient(ctx, bsdkv3.AppkeyPcr,
		bsdkv3.WithClientLogger(logger),
	)
	if err != nil {
		logger.Error("创建客户端失败！")
		logger.Error("错误信息：%s", err.Error())
		return
	}

	// 执行登录并自动处理验证码
	ret, err := client.Login(ctx, user)
	if err != nil {
		logger.Error("登录失败！")
		logger.Error("错误信息：%s", err.Error())
		return
	}

	logger.Info("登录成功！")
	logger.Info("Uid: %v, AccessKey: %s", ret.Uid, ret.AccessKey)
}

func main() {
	fmt.Println("============== test bsdkv3r ==============")
	testLoginR()
}
