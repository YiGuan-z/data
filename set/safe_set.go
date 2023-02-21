package set

import "sync"

type safeSet struct {
	//用于存储数据
	data map[any]struct{}
	//用于统计元素的个数
	size int
	//读写锁
	m sync.RWMutex
}

func (s *safeSet) Add(a any) {
	s.m.Lock()
	defer s.m.Unlock()
	if _, ok := s.data[a]; !ok {
		s.data[a] = struct{}{}
		s.size++
	}
}

func (s *safeSet) Adds(a ...any) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, val := range a {
		if _, ok := s.data[val]; !ok {
			s.data[val] = struct{}{}
			s.size++
		}
	}
}

func (s *safeSet) Remove(a any) (any, bool) {
	s.m.Lock()
	val, ok := s.data[a]
	if ok {
		delete(s.data, a)
		s.size--
	} else {
		return nil, ok
	}
	s.m.Unlock()
	return val, ok
}

func (s *safeSet) Range(f func(any)) {
	s.RangeR(f)
}

func (s *safeSet) Filter(f func(any) bool) Set {
	if f == nil {
		panic("没有定义过滤器方法")
	}
	s.RangeW(func(val any) {
		ok := f(val)
		if !ok {
			delete(s.data, val)
			s.size--
		}
	})
	//修正长度
	if s.size != len(s.data) {
		s.size = len(s.data)
	}
	return s
}

//func (s *safeSet) Find(f func(any) bool) (ret *interface{}) {
//	if f == nil {
//		panic("没有定义查找方法")
//	}
//	s.RangeR(func(val any) {
//		ok := f(val)
//		if ok {
//			ret = &val
//			return
//		}
//	})
//	return
//}

func (s *safeSet) Clear() {
	s.m.Lock()
	defer s.m.Unlock()
	//在读取大小的时候也应该同步,外部上锁就行
	if s.Size() > 0 {
		s.rangeRW(func(val any) {
			delete(s.data, val)
			s.size--
		}, false, false)
	}
}

func (s *safeSet) Size() int {
	return s.size
}

func (s *safeSet) IsEmpty() bool {
	return s.size == 0
}

func (s *safeSet) Contains(a any) bool {
	s.m.RLock()
	defer s.m.RUnlock()
	_, ok := s.data[a]
	return ok
}

func (s *safeSet) ContainsAll(anies []any) bool {
	s.m.RLock()
	defer s.m.RUnlock()
	r := true
	for _, val := range anies {
		_, ok := s.data[val]
		r = r && ok
	}
	return r
}

func (s *safeSet) AddAll(anies []any) bool {
	s.m.Lock()
	defer s.m.Unlock()
	for _, val := range anies {
		s.Add(val)
	}
	return true
}

func (s *safeSet) RetainAll(anies []any) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, val := range anies {
		if value, ok := s.data[val]; !ok {
			s.Remove(value)
		}
	}
}

func (s *safeSet) RemoveAll(anies []any) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, val := range anies {
		if value, ok := s.data[val]; ok {
			s.Remove(value)
		}
	}
}

func (s *safeSet) ToSlice() []any {
	ret := make([]any, s.Size())
	s.RangeR(func(val any) {
		ret = append(ret, val)
	})
	defer s.Clear()
	return ret
}

// rangeRW 传递一个方法，后两个参数决定是否加锁或者加什么类型的锁
func (s *safeSet) rangeRW(f func(val any), r, w bool) {
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
		for k := range s.data {
			f(k)
		}
	}
}

// RangeR 读锁
func (s *safeSet) RangeR(f func(val any)) {
	s.rangeRW(f, true, false)
}

// RangeW 写锁
func (s *safeSet) RangeW(f func(val any)) {
	s.rangeRW(f, false, true)
}
