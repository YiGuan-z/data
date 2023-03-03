package stream

import (
	"github.com/YiGuan-z/data/set"
	"sync/atomic"
)

type Stream interface {
	// Find 过滤Stream中的数据并返回它第一个判定为真的数据
	Find(func(val any) bool) any
	// ToMap 将数据映射为Map
	ToMap(func(val any) (key, value any)) map[any]any
	// ToSet 将Stream转化为Set集合
	ToSet(func(val any) any) set.Set
	// ToArray 将内容转化为数组
	ToArray() []any
	// GroupBy 分类函数 第一个切片为条件判定为true，第二个切片为条件判定为false
	GroupBy(func(any) bool) (yes, no []any)
	// Count 统计有多少元素符合方法约定
	Count(func(val any) bool) int
	// Chanel 转化为管道进行操作
	Chanel() <-chan any
	// Range 用于循环元素
	Range(func(val any))
	// Size 返回元素的个数
	Size() int
	// Filter 过滤Stream中的数据并返回过滤后的Stream,方法返回true保留数据，false则丢弃数据
	Filter(func(val any) bool) Stream
	// Map 将Stream中的数据映射成另一种数据
	Map(func(val any) any) Stream
	// Head 取头部几个元素并返回流
	Head(int) Stream
	// Tail 取尾部几个元素并返回流
	Tail(int) Stream
	// Skip 跳过元素
	Skip(int) Stream
	// Limit 对无限流进行截断
	Limit(int) Stream
	// Distinct 去除重复元素
	Distinct() Stream
}

// NewChanStream 创建管道流
func NewChanStream(data []any) Stream {
	send, lenght := newChanel(data)
	size := new(int)
	*size = lenght
	return &ChanStream{
		p:    send,
		size: size,
	}
}

//	newChanStreamOfChan ChanStream 的私有构造，
//
// 使用它可快速创建下一个流对象进行操作，由于只需要传递一个管道，所以不会有太大的性能损失，
// 数据仅在管道内部进行流动，进行每一个用户定义的数据操作
// 这个管道将会发送多少数据必须标注出来 因为Head Tail Skip 这三个方法需要使用它
func newChanStreamOfChan(data <-chan any, size int, infinite bool) Stream {
	return &ChanStream{
		p:        data,
		size:     &size,
		infinite: infinite,
	}
}

// NewArrayStream 创建数组流
func NewArrayStream(data []any) Stream {
	return &ArrayStream{data: data, size: len(data)}
}

// Generate 生成一个无限流
func Generate(f func() any) Stream {
	if f == nil {
		panic("生成器方法未定义")
	}
	g := make(chan any)
	//创建一个int指针，并将它传递给结构体以便于对生成元素个数的统计
	count := new(int)
	//这个原子变量是用来控制生成器的停止方法，后续可能会有其它方案，现在暂时使用原子变量
	stop := atomic.Bool{}
	stop.Store(true)
	go func() {
		//这个go程会通过原子变量判断是否需要停止生成器方法，通过已生成的流对象的Limit()方法进行截断
		for stop.Load() {
			g <- f()
			*count++
		}
		close(g)
	}()
	return &ChanStream{
		p:            g,
		size:         count,
		infinite:     true,
		stopGenerate: &stop,
	}
}

// Iteration seed 是一个生成的种子，hasNext判断是否继续生成，next是下一个元素的生成具体方法
func Iteration(seed any, hasNext func(any) bool, next func(any) any) Stream {
	checkFunc(hasNext)
	checkFunc(next)

	g := make(chan any)
	count := new(int)
	//启动一个go程，专门负责对元素的生成
	//该闭包捕获了5个变量 hasNext seed next g count，
	//这五个变量尽量不要在其它环境对其修改，如果有move关键字就好了🤔
	go func() {
		//判断是否需要迭代
		for hasNext(seed) {
			//将当前元素传给迭代方法进行迭代
			seed = next(seed)
			//交给管道
			g <- seed
			//对生成的元素+1
			*count++
		}
		close(g)
	}()
	return &ChanStream{
		p:        g,
		size:     count,
		infinite: false,
	}
}
