package set

type Set interface {
	// Add 添加一个元素
	Add(t any)
	// Adds 添加多个元素
	Adds(t ...any)
	// Remove 移除一个元素
	Remove(t any) (any, bool)
	// Range 对元素进行操作
	Range(func(val any))
	// Filter 根据条件删除对应元素 表达式判定为false将删除该元素
	Filter(f func(val any) bool) Set
	// Find 查找一个对象并返回它的指针
	Find(f func(val any) bool) (ret *interface{})
	// Clear 清空元素
	Clear()
	// Size 元素的个数
	Size() int
	// IsEmpty 判断Set集合是否为空
	IsEmpty() bool
	// Contains 判断Set集合是否包含该元素
	Contains(val any) bool
	// ContainsAll 判断Set集合是否包含该切片中的所有元素
	ContainsAll([]any) bool
	// AddAll 添加切片内所有元素到集合中
	AddAll([]any) bool
	// RetainAll 只保留该切片中的所有元素
	RetainAll([]any)
	// RemoveAll 删除在该切片中的所有元素
	RemoveAll([]any)
}
