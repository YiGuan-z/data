package stream

import "sync/atomic"

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

// NewArrayStream 创建数组流
func NewArrayStream(data []any) Stream {
	return &ArrayStream{data: data, size: len(data)}
}

// newChanStreamOfChan ChanStream 的构造，必须携带发送内容的长度，以便应对使用切片、map的时候造成性能损失。
// 并且Head Tail Skip 这三个方法需要使用它
func newChanStreamOfChan(data <-chan any, size int, infinite bool) Stream {
	return &ChanStream{
		p:        data,
		size:     &size,
		infinite: infinite,
	}
}

// newChanel 根据切片创建一个管道并返回大小
func newChanel(data []any) (send <-chan any, lenght int) {
	ret := make(chan any)
	go func() {
		for _, v := range data {
			ret <- v
		}
		close(ret)
	}()
	send = ret
	lenght = len(data)
	return
}

// Generate 生成一个无限流
func Generate(f func() any) Stream {
	if f == nil {
		panic("生成器方法未定义")
	}
	g := make(chan any)
	count := new(int)
	stop := atomic.Bool{}
	stop.Store(true)
	go func(c chan<- any) {
		for stop.Load() {
			c <- f()
			*count++
		}
		close(g)
	}(g)
	ret := &ChanStream{
		p:            g,
		size:         count,
		infinite:     true,
		stopGenerate: &stop,
	}
	return ret
}

func Iteration(seed any, hasNext func(any) bool, next func(any) any) Stream {
	checkFunc(hasNext)
	checkFunc(next)
	g := make(chan any)
	count := new(int)
	go func(c <-chan any) {
		//判断是否需要迭代
		for hasNext(seed) {
			//将当前元素传给迭代方法进行迭代
			seed = next(seed)
			//交给管道
			g <- seed
			*count++
		}
		close(g)
	}(g)
	ret := &ChanStream{
		p:        g,
		size:     count,
		infinite: false,
	}
	return ret
}

func checkFunc(f any) {
	if f == nil {
		panic("方法未定义")
	}
}
