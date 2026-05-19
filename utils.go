package bsdkv3

import "encoding/json"

// 函数式选项模式（包内使用，调用方只需传入 With* 返回的 option 值）

type option[T any] interface {
	apply(*T)
}

// 2. 定义内部的泛型闭包类型
type optionFunc[T any] func(*T)

// 3. 实现接口
func (f optionFunc[T]) apply(i *T) {
	f(i)
}

// optionalValue 区分 JSON 响应中「字段未出现」与「字段出现但为空/null」。
type optionalValue[T any] struct {
	Valid bool // 标记该字段是否在响应中出现过
	Value T    // 实际的值
}

// 实现 json.Unmarshaler 接口
func (o *optionalValue[T]) UnmarshalJSON(data []byte) error {
	// 只要这个方法被调用了，说明 JSON 中绝对传了这个字段！
	o.Valid = true

	// 如果服务器传的是明确的 null，可以选择保留默认零值，或者加一个 isNull 标记
	if string(data) == "null" {
		return nil
	}

	// 解析实际的数据到 Value 中
	return json.Unmarshal(data, &o.Value)
}
