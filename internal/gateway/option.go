package gateway

// 1. 定义一个通用的 Option 接口
type Option[T any] interface {
	apply(*T)
}

// 2. 定义内部的泛型闭包类型
type optionFunc[T any] func(*T)

// 3. 实现接口
//
//nolint:unused // 泛型类型的方法在满足接口时可能会被 unused 误报
func (f optionFunc[T]) apply(i *T) {
	f(i)
}
