package stream

import (
	"github.com/YiGuan-z/data/set"
)

type ArrayStream struct {
	data []any
	size int
}

func (a *ArrayStream) Filter(f func(val any) bool) Stream {
	ret := make([]any, a.size)
	a.Range(func(val any) {
		if f(val) {
			ret = append(ret, val)
		}
	})
	return NewArrayStream(ret)
}

func (a *ArrayStream) Find(f func(val any) bool) (r any) {
	for _, value := range a.data {
		if f(value) {
			r = value
			break
		}
	}
	return
}

func (a *ArrayStream) ToMap(f func(val any) (key any, value any)) map[any]any {
	ret := make(map[any]any, a.size)
	a.Range(func(val any) {
		key, value := f(val)
		ret[key] = value
	})
	return ret
}

func (a *ArrayStream) ToSet(f func(val any) any) set.Set {
	ret := set.NewUnsafeSetOfLen(a.size)
	a.Range(func(val any) {
		entry := f(val)
		ret.Add(entry)
	})
	return ret
}

func (a *ArrayStream) Map(f func(val any) any) Stream {
	for i, val := range a.data {
		a.data[i] = f(val)
	}
	return a
}

func (a *ArrayStream) ToArray() []any {
	arr := make([]any, a.size)
	a.Range(func(val any) {
		arr = append(arr, val)
	})
	return arr
}

func (a *ArrayStream) Count(f func(val any) bool) int {
	count := 0
	a.Range(func(val any) {
		if f(val) {
			count++
		}
	})
	return count
}

func (a *ArrayStream) Chanel() <-chan any {
	ch := make(chan any)
	go func() {
		defer close(ch)
		for _, value := range a.data {
			ch <- value
		}
		return
	}()
	return ch
}

func (a *ArrayStream) Head(i int) Stream {
	s := make([]any, i)
	c := a.data[:i]
	copy(c, s)
	return NewArrayStream(s)
}

func (a *ArrayStream) Tail(i int) Stream {
	s := make([]any, i)
	c := a.data[len(a.data)-i:]
	copy(c, s)
	return NewArrayStream(s)
}

func (a *ArrayStream) Skip(i int) Stream {
	s := make([]any, i)
	c := a.data[i:]
	copy(c, s)
	return NewArrayStream(s)
}

func (a *ArrayStream) Range(f func(val any)) {
	for _, val := range a.data {
		f(val)
	}
}

func (a *ArrayStream) Size() int {
	return a.size
}

func (a *ArrayStream) Distinct() Stream {
	store := set.NewUnsafeSetOfLen(a.size)
	a.Range(func(val any) {
		store.Add(val)
	})
	return NewChanStream(store.ToSlice())
}

func (a *ArrayStream) Limit(i int) Stream {
	return a.Head(i)
}

func (a *ArrayStream) GroupBy(f func(any) bool) (yes, no []any) {
	yes = make([]any, a.size)
	no = make([]any, a.size)
	a.Range(func(val any) {
		if f(val) {
			yes = append(yes, val)
		} else {
			no = append(no, val)
		}
	})
	return
}
