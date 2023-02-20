package set

import (
	"sync"
)

type AnySet struct {
	data  map[any]struct{}
	count int
	m     sync.RWMutex
}

//type Entry struct {
//	p unsafe.Pointer
//}

func NewSet(size int) Set {
	return &AnySet{
		data: make(map[any]struct{}, size),
	}
}

// Add 添加一个元素
func (s *AnySet) Add(t any) {
	s.m.Lock()
	defer s.m.Unlock()
	if _, ok := s.data[t]; !ok {
		s.data[t] = struct{}{}
		s.count++
	}
}

// Adds 添加多个元素
func (s *AnySet) Adds(t ...any) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, val := range t {
		if _, ok := s.data[val]; !ok {
			s.data[val] = struct{}{}
			s.count++
		}
	}
}

// Remove 移除一个元素
func (s *AnySet) Remove(t any) (any, bool) {
	s.m.Lock()
	val, ok := s.data[t]
	if ok {
		delete(s.data, t)
		s.count--
	} else {
		return nil, ok
	}
	s.m.Unlock()
	return val, ok
}

// Range 对元素进行操作
func (s *AnySet) Range(f func(val any)) {
	//对接下来的操作上读锁
	s.RangeR(f)
}

// rangeRW 传递一个用于方法，后两个参数决定是否加锁或者加什么类型的锁
func (s *AnySet) rangeRW(f func(val any), r, w bool) {
	if r && w {
		panic("不能同时开启读写锁")
	}
	if r {
		s.m.RLock()
		defer s.m.RUnlock()
	}
	if w {
		s.m.Lock()
		defer s.m.Unlock()
	}

	if f == nil {
		panic("方法未定义")
	}
	if s.Size() > 0 {
		for k, _ := range s.data {
			f(k)
		}
	}
}

// RangeR 读锁
func (s *AnySet) RangeR(f func(val any)) {
	s.rangeRW(f, true, false)
}

// RangeW 写锁
func (s *AnySet) RangeW(f func(val any)) {
	s.rangeRW(f, false, true)
}

// Clear 清空元素
func (s *AnySet) Clear() {
	s.m.Lock()
	defer s.m.Unlock()
	//在读取大小的时候也应该同步,外部上锁就行
	if s.Size() > 0 {
		s.rangeRW(func(val any) {
			delete(s.data, val)
			s.count--
		}, false, false)
	}
}

// Filter 根据条件删除对应元素条件判定为false的时候删除该元素，返回的是原本对象
func (s *AnySet) Filter(f func(val any) bool) Set {
	if f == nil {
		panic("没有定义过滤器方法")
	}
	s.RangeW(func(val any) {
		ok := f(val)
		if !ok {
			delete(s.data, val)
			s.count--
		}
	})
	if s.count != len(s.data) {
		s.count = len(s.data)
	}
	return s
}

// CloneFilter 根据条件删除对应元素，条件判定为false的时候删除该元素，返回的是一个新对象。
//func (s AnySet) CloneFilter(f func(val any) bool) Set {
//	if f == nil {
//		panic("没有定义过滤器方法")
//	}
//	s.RangeW(func(val any) {
//		ok := f(val)
//		if !ok {
//			delete(s.data, val)
//		}
//	})
//	return &s
//}

// Find 通过关键特征查找一个对象并返回它的指针
func (s *AnySet) Find(f func(val any) bool) (ret *interface{}) {
	if f == nil {
		panic("没有定义查找方法")
	}
	s.RangeR(func(val any) {
		ok := f(val)
		if ok {
			ret = &val
			return
		}
	})
	ret = nil
	return
}

func (s *AnySet) IsEmpty() bool {
	return s.count == 0
}

func (s *AnySet) Contains(val any) bool {
	s.m.RLock()
	defer s.m.RUnlock()
	_, ok := s.data[val]
	return ok
}

func (s *AnySet) RetainAll(anies []any) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, val := range anies {
		if value, ok := s.data[val]; !ok {
			s.Remove(value)
		}
	}
}

func (s *AnySet) RemoveAll(anies []any) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, val := range anies {
		if value, ok := s.data[val]; ok {
			s.Remove(value)
		}
	}
}

func (s *AnySet) ContainsAll(anies []any) bool {
	s.m.RLock()
	defer s.m.RUnlock()
	r := true
	for _, val := range anies {
		_, ok := s.data[val]
		r = r && ok
	}
	return r
}

func (s *AnySet) AddAll(anies []any) bool {
	s.m.Lock()
	defer s.m.Unlock()
	for _, val := range anies {
		s.Add(val)
	}
	return true
}

func (s *AnySet) Size() int {
	return s.count
}

func (s *AnySet) ToSlice() []any {
	ret := make([]any, s.Size())
	s.RangeR(func(val any) {
		ret = append(ret, val)
	})
	defer s.Clear()
	return ret
}
