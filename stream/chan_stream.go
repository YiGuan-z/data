package stream

import (
	"github.com/YiGuan-z/data/set"
)

type ChanStream struct {
	p    <-chan any
	size int
}

func NewChanStream(data []any) Stream {
	send, lenght := newChanel(data)
	return &ChanStream{
		p:    send,
		size: lenght,
	}
}

// newChanStreamOfChan ChanStream 的构造，必须携带发送内容的长度，以便应对初始化切片、map的时候造成性能损失。
// 并且Head Tail Skip 这三个方法需要使用它
func newChanStreamOfChan(data <-chan any, size int) Stream {
	return &ChanStream{
		p:    data,
		size: size,
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

func (c *ChanStream) Find(f func(val any) bool) any {
	var ret any
	c.DefaultRange(func(val any) {
		if f(val) {
			ret = val
		}
	})
	return ret
}

func (c *ChanStream) ToMap(f func(val any) (key any, value any)) map[any]any {
	retMap := make(map[any]any, c.size)
	c.DefaultRange(func(val any) {
		key, value := f(val)
		retMap[key] = value
	})
	return retMap
}

func (c *ChanStream) ToSet(f func(val any) any) set.Set {
	ret := set.NewSet(c.size)
	c.DefaultRange(func(val any) {
		ret.Add(f(val))
	})
	return ret
}

func (c *ChanStream) ToArray() []any {
	arr := make([]any, c.size)
	c.DefaultRange(func(val any) {
		arr = append(arr, val)
	})
	return arr
}

func (c *ChanStream) Count(f func(val any) bool) int {
	count := new(int)
	c.DefaultRange(func(val any) {
		if f(val) {
			*count++
		}
	})
	return *count
}

func (c *ChanStream) Chanel() <-chan any {
	ch := make(chan any)
	go func() {
		for value := range c.p {
			ch <- value
		}
		close(ch)
	}()
	return ch
}

func (c *ChanStream) Filter(f func(val any) bool) Stream {
	ch := make(chan any)
	go c.RangeFunc(func(val any) {
		if f(val) {
			ch <- val
		} else {
			//被舍弃掉的元素
			c.size--
		}
	}, func() {
		close(ch)
	}, 0)
	return newChanStreamOfChan(ch, c.size)
}

func (c *ChanStream) Map(f func(val any) any) Stream {
	arr := make([]any, c.size)
	c.DefaultRange(func(val any) {
		obj := f(val)
		arr = append(arr, obj)
	})
	return NewChanStream(arr)
}

func (c *ChanStream) Head(i int) Stream {
	count := 0
	r := make(chan any)
	go c.RangeFunc(func(val any) {
		if count < i {
			r <- val
			//因为只取头部几个元素，所以只在if语句里对count自增，获取完毕就不需要了
			count++
		}
	}, func() {
		close(r)
	}, 0)
	return newChanStreamOfChan(r, i)
}

func (c *ChanStream) Tail(i int) Stream {
	count := 0
	r := make(chan any)
	size := c.size - i
	go c.RangeFunc(func(val any) {
		if count >= size {
			r <- val
		}
		count++
	}, func() {
		close(r)
	}, 0)
	return newChanStreamOfChan(r, i)
}

func (c *ChanStream) Skip(i int) Stream {
	r := make(chan any)
	go c.RangeFunc(func(val any) {
		r <- val
	}, func() {
		close(r)
	}, i)
	return newChanStreamOfChan(r, c.size-i)
}

func (c *ChanStream) Range(f func(val any)) {
	c.DefaultRange(f)
}

func (c *ChanStream) Size() int {
	return c.size
}

// RangeFunc f是循环内的操作函数，end代表循环结束后的收尾操作，offset代表偏移量，丢弃掉一些数据
func (c *ChanStream) RangeFunc(f func(val any), end func(), offset int) {
	count := 0
	for val := range c.p {
		//如果计数器小于偏移量，则跳过
		if count < offset {
			//长度计算器计算被忽略的元素
			c.size--
			count++
			continue
		}
		f(val)
		count++
	}
	end()
}

func (c *ChanStream) DefaultRange(f func(val any)) {
	c.RangeFunc(f, func() {}, 0)
}
