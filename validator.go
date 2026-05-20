package bsdkv3

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"resty.dev/v3"
)

// ValidatorResult 定义验证结果的标准化结构
type ValidatorResult struct {
	Challenge string `json:"challenge"`
	Gt        string `json:"gt"`
	GtUserId  string `json:"gt_user_id"`
	Validate  string `json:"validate"`
}

// Validator 定义验证码校验行为。
type Validator interface {
	Validate(ctx context.Context) (*ValidatorResult, error)
}

// remoteValidator 远程验证服务实现
type remoteValidator struct {
	httpClient   *resty.Client
	logger       Logger
	tryTimes     int
	pollInterval time.Duration
}

type remoteValidatorOptions struct {
	logger       Logger
	tryTimes     int
	pollInterval time.Duration
}

const (
	remoteValidatorEndpoint      = "https://pcrd.tencentbot.top/geetest_renew"
	remoteValidatorCheckEndpoint = "https://pcrd.tencentbot.top/check"

	defaultTryTimes = 5

	// 远程验证服务队列轮询阈值
	maxAcceptableQueueLength = 35
	maxQueueWaitSlots        = 3
	queueSlotWaitSeconds     = 5

	defaultPollIntervalSeconds = 8
)

var (
	defaultPollInterval = defaultPollIntervalSeconds * time.Second
	defaultUserAgent    = "autopcr/1.0.0"
)

type getValidationUUIDResp struct {
	UUID string `json:"uuid"`
}

func withValidatorLogger(logger Logger) option[remoteValidatorOptions] {
	return optionFunc[remoteValidatorOptions](func(o *remoteValidatorOptions) {
		o.logger = logger
	})
}

// 保留备用
// func withValidatorTryTimes(times int) option[remoteValidatorOptions] {
// 	return optionFunc[remoteValidatorOptions](func(o *remoteValidatorOptions) {
// 		o.tryTimes = times
// 	})
// }

// func withValidatorPollInterval(interval time.Duration) option[remoteValidatorOptions] {
// 	return optionFunc[remoteValidatorOptions](func(o *remoteValidatorOptions) {
// 		o.pollInterval = interval
// 	})
// }

func newRemoteValidator(httpClient *resty.Client, opts ...option[remoteValidatorOptions]) Validator {
	// 默认配置
	options := remoteValidatorOptions{
		logger:       discardLogger,
		tryTimes:     defaultTryTimes,
		pollInterval: defaultPollInterval,
	}

	// 应用函数式选项
	for _, o := range opts {
		o.apply(&options)
	}

	if options.logger == nil {
		options.logger = discardLogger
	}

	return &remoteValidator{
		httpClient: httpClient.
			SetHeader("User-Agent", defaultUserAgent).
			SetHeader("Content-Type", "application/json"),
		logger:       options.logger,
		tryTimes:     options.tryTimes,
		pollInterval: options.pollInterval,
	}
}

func (r *remoteValidator) getValidationUUID(ctx context.Context) (string, error) {
	r.logger.Debug("开始获取验证UUID")
	resp, err := r.httpClient.R().
		SetContext(ctx).
		Get(remoteValidatorEndpoint)
	if err != nil {
		r.logger.Error("请求验证UUID时发生错误: %v", err)
		return "", fmt.Errorf("网络请求验证UUID失败: %w", err)
	}

	if resp.IsError() {
		r.logger.Error("获取验证UUID请求失败。状态码: %d", resp.StatusCode())
		return "", fmt.Errorf("验证UUID请求返回错误状态码: %d", resp.StatusCode())
	}

	var result getValidationUUIDResp
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		r.logger.Error("解析验证UUID响应失败: %v", err)
		return "", fmt.Errorf("解析验证UUID响应JSON失败: %w", err)
	}

	if result.UUID == "" {
		r.logger.Error("响应中不包含有效的UUID")
		return "", fmt.Errorf("响应中不包含有效的UUID字段或格式错误")
	}

	return result.UUID, nil
}

func (r *remoteValidator) pollValidationStatus(ctx context.Context, uuid string) (*ValidatorResult, error) {
	maxAttempts := r.tryTimes
	checkEndpoint, err := url.JoinPath(remoteValidatorCheckEndpoint, uuid)
	if err != nil {
		r.logger.Error("构建检查端点URL失败: %v", err)
		return nil, fmt.Errorf("构建检查端点URL失败: %w", err)
	}
	r.logger.Debug("开始轮询验证状态，最大尝试次数: %d", maxAttempts)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result := map[string]any{}
		r.logger.Debug("验证状态检查 #%d/%d", attempt, maxAttempts)
		resp, err := r.httpClient.R().
			SetContext(ctx).
			SetResult(&result).
			Get(checkEndpoint)
		if err != nil {
			r.logger.Error("请求验证状态时发生错误: %v", err)
			return nil, fmt.Errorf("请求验证状态失败: %w", err)
		}

		if resp.IsError() {
			r.logger.Error("检查请求失败，状态码: %d", resp.StatusCode())
			return nil, fmt.Errorf("check request failed with status: %d", resp.StatusCode())
		}

		if queueNum, ok := result["queue_num"].(float64); ok {
			if queueNum >= maxAcceptableQueueLength {
				r.logger.Error("验证码队列过长")
				return nil, fmt.Errorf("captcha failed: queue too long")
			}
			q := min(int(queueNum), maxQueueWaitSlots)
			sleepDuration := time.Duration(q) * queueSlotWaitSeconds * time.Second
			if err := sleepContext(ctx, sleepDuration); err != nil {
				r.logger.Error("等待验证状态时发生错误: %v", err)
				return nil, fmt.Errorf("等待验证状态失败: %w", err)
			}
			continue
		}

		info, ok := result["info"]
		if !ok {
			r.logger.Error("响应中缺少info字段")
			return nil, fmt.Errorf("captcha failed: missing info field")
		}

		switch v := info.(type) {
		case string:
			switch v {
			case "fail", "url invalid":
				r.logger.Error("验证失败: %s", v)
				return nil, fmt.Errorf("captcha failed: %s", v)
			case "in running":
				if err := sleepContext(ctx, r.pollInterval); err != nil {
					r.logger.Error("等待验证状态时发生错误: %v", err)
					return nil, fmt.Errorf("等待验证状态失败: %w", err)
				}
				continue
			default:
				r.logger.Warn("未知的验证状态: %s，继续等待...", v)
				if err := sleepContext(ctx, r.pollInterval); err != nil {
					r.logger.Error("等待验证状态时发生错误: %v", err)
					return nil, fmt.Errorf("等待验证状态失败: %w", err)
				}
				continue
			}
		case map[string]any:
			challenge, _ := v["challenge"].(string)
			gt, _ := v["gt"].(string)
			gtUserId, _ := v["gt_user_id"].(string)
			validate, _ := v["validate"].(string)

			if challenge == "" || gt == "" || validate == "" {
				r.logger.Error("验证响应缺少必填字段")
				return nil, fmt.Errorf("validation response missing required fields")
			}
			r.logger.Info("验证成功完成")
			return &ValidatorResult{
				Challenge: challenge,
				Gt:        gt,
				GtUserId:  gtUserId,
				Validate:  validate,
			}, nil
		}
	}

	return nil, fmt.Errorf("captcha failed: max attempts reached")
}

func (r *remoteValidator) Validate(ctx context.Context) (*ValidatorResult, error) {
	r.logger.Info("使用远程验证服务 (Powered by lulu)")

	uuid, err := r.getValidationUUID(ctx)
	if err != nil {
		r.logger.Error("获取验证UUID失败: %v", err)
		return nil, err
	}
	r.logger.Debug("获取到验证UUID: %s", uuid)

	return r.pollValidationStatus(ctx, uuid)
}

func sleepContext(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
