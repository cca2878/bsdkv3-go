package bsdkv3

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"resty.dev/v3"
)

func (c *Client) calcSign(values url.Values) (string, error) {
	keys := make([]string, 0, len(values))
	keys = keys[:len(values)]
	i := 0
	for k := range values {
		keys[i] = k
		i++
	}
	keys = keys[:i]
	sort.Strings(keys)

	h := md5.New()
	for _, k := range keys {
		io.WriteString(h, values.Get(k))
	}
	io.WriteString(h, c.appKey)

	return hex.EncodeToString(h.Sum(nil)), nil
}

func (c *Client) preReq(ctx context.Context, reqBody any) (*resty.Request, error) {
	values, err := c.formEncoder.Encode(reqBody)
	if err != nil {
		return nil, fmt.Errorf("请求体编码失败: %w", err)
	}

	if values.Get("timestamp") == "" {
		values.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	}

	sign, err := c.calcSign(values)
	if err != nil {
		return nil, fmt.Errorf("计算签名失败: %w", err)
	}

	values.Set("sign", sign)

	return c.httpClient.R().
		SetContext(ctx).
		SetFormDataFromValues(values), nil
}

func (c *Client) execReq(ctx context.Context, method, path string, reqBody any, respBody any) (*resty.Response, error) {
	// requ, err := c.preReq(ctx, reqBody)
	// if err != nil {
	// 	c.logger.Error("准备请求失败: %v", err)
	// 	return nil, err
	// }

	// url, err := req.URL(c.hosts)
	// if err != nil {
	// 	c.logger.Error("获取URL失败: %v", err)
	// 	return nil, err
	// }

	// c.logger.Debug("发送请求: %s", url.String())

	// resp, err := requ.SetResult(result).Execute(req.Method(), url.String())
	// if err != nil {
	// 	c.logger.Error("请求发送失败: %v", err)
	// 	return resp, err
	// }

	// c.logger.Debug("收到响应: 状态码=%d, 内容长度=%d", resp.StatusCode(), len(resp.Body()))

	// if resp.StatusCode() != http.StatusOK {
	// 	c.logger.Error("请求失败，状态码: %d", resp.StatusCode())
	// 	return resp, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode())
	// }

	// return resp, nil
	var lastErr error

	for i := 0; i < c.clientConf.TryTimes; i++ {
		var currentHost string
		var targetURL *url.URL
		var err error

		if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
			targetURL, err = url.Parse(path)
		} else {
			currentHost = c.hostMgr.getCurrentHost()
			targetURL, err = url.Parse(currentHost + path)
		}

		if err != nil {
			return nil, fmt.Errorf("解析URL失败: %w", err)
		}

		// preReq 内部去计算签名、处理表单
		req, err := c.preReq(ctx, reqBody)
		if err != nil {
			return nil, err
		}

		resp, err := req.SetResult(respBody).Execute(method, targetURL.String())

		if err != nil || resp.StatusCode() >= 500 {
			if currentHost != "" {
				c.hostMgr.markFailed(currentHost)
			}
			if err != nil {
				lastErr = fmt.Errorf("请求发送失败: %w", err)
			} else {
				lastErr = fmt.Errorf("请求失败，状态码: %d", resp.StatusCode())
			}
			continue
		}

		return resp, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("未执行任何请求或无可用错误信息")
	}
	return nil, fmt.Errorf("请求失败: %w", lastErr)
}

func execAPI[Req request, Resp any](ctx context.Context, c *Client, endpoint endpoint[Req, Resp], reqBody Req) (*apiResponse[Resp], error) {
	if err := reqBody.validate(); err != nil {
		return nil, fmt.Errorf("请求参数校验失败: %w", err)
	}

	respBody := new(Resp)
	restyResp, err := c.execReq(ctx, endpoint.Method, endpoint.Path, reqBody, respBody)
	if err != nil {
		return nil, err
	}
	return &apiResponse[Resp]{
		Body: *respBody,
		Raw:  restyResp,
	}, nil
}
