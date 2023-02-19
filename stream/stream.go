package stream

import (
	"github.com/YiGuan-z/data/set"
)

type Stream interface {
	// Filter 过滤Stream中的数据并返回过滤后的Stream
	Filter(func(val any) bool) Stream
	// Find 过滤Stream中的数据并返回它第一个判定为真的数据
	Find(func(val any) bool) any
	// ToMap 将数据映射为Map
	ToMap(func(val any) (key, value any)) map[any]any
	// ToSet 将Stream转化为Set集合
	ToSet(func(val any) any) set.Set
	// ToArray 将内容转化为数组
	ToArray() []any
	// Map 将Stream中的数据映射成另一种数据
	Map(func(val any) any) Stream
	// Count 统计有多少元素符合方法约定
	Count(func(val any) bool) int
	// Chanel 转化为管道进行操作
	Chanel() <-chan any
	// Head 取头部几个元素并返回流
	Head(int) Stream
	// Tail 取尾部几个元素并返回流
	Tail(int) Stream
	// Skip 跳过元素
	Skip(int) Stream
	// Range 用于循环元素
	Range(func(val any))
}