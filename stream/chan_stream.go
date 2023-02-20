package stream

import (
	"github.com/YiGuan-z/data/set"
	"sync/atomic"
)

type ChanStream struct {
	p            <-chan any
	size         int
	infinite     bool
	stopGenerate *atomic.Bool
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
	count := 0
	c.DefaultRange(func(val any) {
		if f(val) {
			count++
		}
	})
	return count
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
			//判断一下是否是无限流
			if !c.infinite {
				c.size--
			}
		}
	}, func() {
		close(ch)
	}, 0)
	return newChanStreamOfChan(ch, c.size, c.infinite)
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
	return newChanStreamOfChan(r, i, c.infinite)
}

func (c *ChanStream) Tail(i int) Stream {
	r := make(chan any)
	size := c.size - i
	go c.RangeFunc(func(val any) {
		r <- val
	}, func() {
		close(r)
	}, size)
	return newChanStreamOfChan(r, i, c.infinite)
}

func (c *ChanStream) Skip(i int) Stream {
	r := make(chan any)
	go c.RangeFunc(func(val any) {
		r <- val
	}, func() {
		close(r)
	}, i)
	return newChanStreamOfChan(r, c.size-i, c.infinite)
}

func (c *ChanStream) Range(f func(val any)) {
	c.DefaultRange(f)
}

func (c *ChanStream) Size() int {
	return c.size
}

func (c *ChanStream) Distinct() Stream {
	store := set.NewSet(c.size)
	c.Range(func(val any) {
		store.Add(val)
	})
	return NewChanStream(store.ToSlice())
}

func (c *ChanStream) Limit(i int) Stream {
	//如果是无限流就进行截断操作
	//不是无限流就转到Skip方法
	if c.infinite {
		count := atomic.Int32{}
		retCh := make(chan any)
		go func() {
			for count.Load() > int32(i) {
				retCh <- c.p
				count.Add(1)
			}
			c.stopGenerate.Store(false)
			close(retCh)
		}()
		return newChanStreamOfChan(retCh, int(count.Load()), false)
	} else {
		return c.Head(i)
	}
}

// RangeFunc f是循环内的操作函数，end代表循环结束后的收尾操作，offset代表偏移量，丢弃掉一些数据
func (c *ChanStream) RangeFunc(f func(val any), end func(), offset int) {
	count := 0
	for val := range c.p {
		//如果计数器小于偏移量，则跳过
		if count < offset {
			//长度计算器计算被忽略的元素
			if !c.infinite {
				c.size--
			}
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
