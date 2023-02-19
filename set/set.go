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
	// Clone 获取该对象的拷贝,没有对set内部对象进行复制，还是引用的原本对象。
	Clone() Set
	// Filter 根据条件删除对应元素
	Filter(f func(val any) bool) Set
	// Find 查找一个对象并返回它的指针
	Find(f func(val any) bool) (ret *interface{})
	// Clear 清空元素
	Clear()
	// Size 元素的个数
	Size() int
}
