package base

import "encoding/json"

// OptionalValue 区分 JSON 响应中「字段未出现」与「字段出现但为空/null」。
type OptionalValue[T any] struct {
	Valid bool // 标记该字段是否在响应中出现过
	Value T    // 实际的值
}

// 实现 json.Unmarshaler 接口
func (o *OptionalValue[T]) UnmarshalJSON(data []byte) error {
	// 只要这个方法被调用了，说明 JSON 中绝对传了这个字段！
	o.Valid = true

	// 如果服务器传的是明确的 null，可以选择保留默认零值，或者加一个 isNull 标记
	if string(data) == "null" {
		return nil
	}

	// 解析实际的数据到 Value 中
	return json.Unmarshal(data, &o.Value)
}
