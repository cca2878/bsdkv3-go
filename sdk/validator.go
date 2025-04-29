package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math"
	"time"

	"bsdkv3/sdk/log"
)

// ValidatorResult 定义验证结果的标准化结构
type ValidatorResult struct {
	Challenge string `json:"challenge"`
	Gt        string `json:"gt"`
	GtUserId  string `json:"gt_user_id"`
	Validate  string `json:"validate"`
}

// Validator 接口定义验证器的行为
type Validator interface {
	Validate() (ValidatorResult, error)
}

// RemoteValidator 远程验证服务实现
type RemoteValidator struct{}

// NewRemoteValidator 创建远程验证器实例
func NewRemoteValidator() Validator {
	return &RemoteValidator{}
}

// Validate 实现 Validator 接口
func (r *RemoteValidator) Validate() (ValidatorResult, error) {
	log.Info("使用远程验证服务 (Powered by lulu)")

	client := resty.New()
	headers := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "autopcr/1.0.0",
	}

	// 获取验证会话的UUID
	uuid, err := getValidationUUID(client, headers)
	if err != nil {
		log.Error("获取验证UUID失败: %v", err)
		return ValidatorResult{}, err
	}
	log.Debug("获取到验证UUID: %s", uuid)

	// 轮询验证状态
	result, err := pollValidationStatus(client, headers, uuid)
	if err != nil {
		return ValidatorResult{}, err
	}

	// 提取必要字段并返回标准化结果
	validatorResult := ValidatorResult{
		Challenge: result["challenge"].(string),
		Gt:        result["gt"].(string),
		GtUserId:  result["gt_user_id"].(string),
		Validate:  result["validate"].(string),
	}

	return validatorResult, nil
}

func getValidationUUID(client *resty.Client, headers map[string]string) (string, error) {
	log.Debug("开始获取验证UUID")
	resp, err := client.R().
		SetHeaders(headers).
		Get("https://pcrd.tencentbot.top/geetest_renew")
	if err != nil {
		log.Error("请求验证UUID时发生错误: %v", err)
		return "", fmt.Errorf("网络请求验证UUID失败: %w", err)
	}

	if resp.IsError() {
		log.Error("获取验证UUID请求失败，状态码: %d", resp.StatusCode())
		return "", fmt.Errorf("验证UUID请求返回错误状态码: %d", resp.StatusCode())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		log.Error("解析验证UUID响应失败: %v", err)
		return "", fmt.Errorf("解析验证UUID响应JSON失败: %w", err)
	}

	uuid, ok := result["uuid"].(string)
	if !ok {
		log.Error("响应中不包含有效的UUID")
		return "", fmt.Errorf("响应中不包含有效的UUID字段或格式错误")
	}

	return uuid, nil
}

// pollValidationStatus 轮询验证服务直到完成
func pollValidationStatus(client *resty.Client, headers map[string]string, uuid string) (map[string]interface{}, error) {
	maxAttempts := 5
	checkEndpoint := fmt.Sprintf("https://pcrd.tencentbot.top/check/%s", uuid)
	log.Debug("开始轮询验证状态，最大尝试次数: %d", maxAttempts)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		log.Debug("验证状态检查 #%d/%d", attempt, maxAttempts)
		// 获取当前验证状态
		resp, err := client.R().
			SetHeaders(headers).
			Get(checkEndpoint)
		if err != nil {
			log.Error("检查验证状态时发生错误: %v", err)
			return nil, err
		}

		if resp.IsError() {
			log.Error("验证状态检查请求失败，状态码: %d", resp.StatusCode())
			return nil, fmt.Errorf("check request failed with status: %d", resp.StatusCode())
		}

		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			log.Error("解析验证状态响应失败: %v", err)
			return nil, err
		}

		log.Debug("验证检查结果 #%d: %v", attempt, result)

		// 处理队列位置
		if queueNum, ok := result["queue_num"].(float64); ok {
			if queueNum >= 35 {
				log.Warn("验证失败: 队列过长 (位置: %.0f)", queueNum)
				return nil, errors.New("captcha failed: queue too long")
			}

			sleepTime := time.Duration(math.Min(queueNum, 3)*10) * time.Second
			log.Debug("UUID %s 处于队列位置 %.0f, 等待 %.0f 秒", uuid, queueNum, sleepTime.Seconds())
			time.Sleep(sleepTime)

			// 长等待时调整尝试计数
			if sleepTime >= 40*time.Second {
				attempt += 2
				log.Debug("等待时间较长，调整尝试计数: %d", attempt)
			}
			continue
		}

		// 处理 info 字段
		info, exists := result["info"]
		if !exists {
			log.Error("验证失败: 响应中缺少 info 字段")
			return nil, errors.New("captcha failed: missing info field")
		}

		// 根据类型和值处理 info
		switch v := info.(type) {
		case string:
			// 处理字符串类型的 info
			switch v {
			case "fail", "url invalid":
				log.Error("验证失败: %s", v)
				return nil, errors.New("captcha failed: " + v)
			case "in running":
				log.Info("验证正在进行中，继续等待...")
				time.Sleep(8 * time.Second)
				continue
			default:
				log.Warn("未知的验证状态: %s", v)
				continue
			}

		case map[string]interface{}:
			// 检查是否包含验证数据和所有必要字段
			if _, hasValidate := v["validate"]; hasValidate {
				// 检查所有必需的字段
				requiredFields := []string{"challenge", "gt", "gt_user_id", "validate"}
				var missingFields []string

				for _, field := range requiredFields {
					if _, exists := v[field]; !exists {
						missingFields = append(missingFields, field)
					}
				}

				if len(missingFields) > 0 {
					log.Error("验证响应缺少必要字段: %v", missingFields)
					return nil, fmt.Errorf("validation response missing required fields: %v", missingFields)
				}

				log.Info("验证成功完成")
				return v, nil
			}
		}
	}

	log.Error("验证失败: 达到最大尝试次数 (%d)", maxAttempts)
	return nil, errors.New("captcha failed: max attempts reached")
}
