package set

type Set interface {
	// Add 添加一个元素
	Add(any)
	// Adds 添加多个元素
	Adds(...any)
	// Remove 移除一个元素
	Remove(any) (any, bool)
	// Range 对元素进行操作
	Range(func(any))
	// Filter 根据条件删除对应元素 表达式判定为false将删除该元素
	Filter(func(any) bool) Set
	// Clear 清空元素
	Clear()
	// Size 元素的个数
	Size() int
	// IsEmpty 判断Set集合是否为空
	IsEmpty() bool
	// Contains 判断Set集合是否包含该元素
	Contains(any) bool
	// ContainsAll 判断Set集合是否包含该切片中的所有元素
	ContainsAll([]any) bool
	// AddAll 添加切片内所有元素到集合中
	AddAll([]any) bool
	// RetainAll 只保留该切片中的所有元素
	RetainAll([]any)
	// RemoveAll 删除在该切片中的所有元素
	RemoveAll([]any)
	// ToSlice 是一个出口方法将元素转化为切片返回并且清除内部元素
	ToSlice() []any
}
