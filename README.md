# bsdkv3-go

> bilibili 游戏 SDK（**external v3 登录协议**）的 Go 客户端——把账号密码换成长期凭据 `access_key`，含风控/验证码处理。

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-CC%20BY--NC--SA%204.0-lightgrey)

只负责 bilibili 登录的**冷启动**这一段：拉取登录 host 列表 + RSA 登录密钥（bootstrap，fail-fast）→
RSA+AES 加密提交账密 → 触发 geetest 风控时经**注入的求解器**过关 → 返回 `(uid, access_key)`。拿到
`access_key` 后即可交给 [go-autopcr-core](https://github.com/cca2878/go-autopcr-core) 等登录游戏。

## 用法

```sh
go get github.com/cca2878/bsdkv3-go
```

```go
import bsdkv3 "github.com/cca2878/bsdkv3-go"

client, err := bsdkv3.NewClient(ctx, bsdkv3.AppkeyPcr) // 隐式 bootstrap，fail-fast
if err != nil { /* ... */ }

acc, err := client.Auth.Login(ctx, bsdkv3.UserInfo{
    Username: "your_account",
    Password: "your_password",
})
// acc.Uid, acc.AccessKey
```

**风控 / 验证码**：账密登录触发 geetest 时需注入求解器（经 `WithClientValidator`）；
[gtrv-go](https://github.com/cca2878/gtrv-go) 的远程求解器实现了 `gtrv.Validator`：

```go
import "github.com/cca2878/gtrv-go"

client, _ := bsdkv3.NewClient(ctx, bsdkv3.AppkeyPcr,
    bsdkv3.WithClientValidator(gtrv.NewRemoteValidator(nil)),
)
```

完整可跑示例见 [`examples/basic`](examples/basic)。

## 选项

| Option | 作用 |
|--------|------|
| `WithClientValidator(v)` | 注入验证码求解器（geetest 风控），实现 `gtrv.Validator` |
| `WithClientGateway(gw)` | 自定义传输网关（默认 HTTP） |
| `WithClientTransport(rt)` | 自定义 `*http.Transport`（代理 / 超时） |
| `WithClientBaseParams(p)` | 覆盖设备 / 基础请求参数 |
| `WithClientRetryTimes(n)` | 登录重试次数 |

## 结构

- 根包 `bsdkv3`：`Client` / `NewClient` / `Auth.Login` + 选项 + 类型（`UserInfo` / `SdkAccount` / `LoginError`）。
- `transport`：网关抽象（`Gateway` / `HTTPGateway`）+ host 管理（多 host 故障转移）。
- `config`：基础请求参数 / 设备信息。
- `internal`：登录服务、RSA/AES 加密、协议编解码（Go internal 规则隐藏）。

## 开发

需要 Go 1.25+。CI（`.github/workflows/ci.yml`）：gofmt / vet / `go mod verify` / build / test（离线单测，
httptest / 占位 host，无真实业务调用）。分支模型 dev/main（开发提 dev、PR 合 main），提交英文
`<type>: <describe>`。

## 许可证

**CC BY-NC-SA 4.0**（署名-非商业性使用-相同方式共享）。仅供个人学习与研究使用。
