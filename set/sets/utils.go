package sets

import (
	"errors"
	"github.com/YiGuan-z/data/set"
	"github.com/YiGuan-z/data/stream"
)

var (
	newStreamError = errors.New("没有选择需要生成的流")
)

// toStream 将Set转化为Stream
func toStream(set set.Set, chanStream, arrayStream bool) (stream.Stream, error) {
	source := make([]any, set.Size())
	set.Range(func(val any) {
		source = append(source, val)
	})
	if arrayStream {
		return stream.NewArrayStream(source), nil
	}
	if chanStream {
		return stream.NewChanStream(source), nil
	}
	return nil, newStreamError
}

// ToArrayStream 将Set转化为基于数组的Stream
func ToArrayStream(set set.Set) (stream.Stream, error) {
	return toStream(set, false, true)
}

// ToChanStream 将Set转化为基于管道和协程的Stream
func ToChanStream(set set.Set) (stream.Stream, error) {
	return toStream(set, true, false)
}

// UnsafeSetToSafeSet 将不安全的Set转化为安全的Set
func UnsafeSetToSafeSet(s set.Set) set.Set {
	data := s.ToSlice()
	return set.NewSafeSetOfSlice(data)
}

// SafeSetToUnsafeSet 将安全的Set转化为不安全的Set
func SafeSetToUnsafeSet(s set.Set) set.Set {
	data := s.ToSlice()
	return set.NewUnsafeSetOfSlice(data)
}

// NewSafeSetOfStream 将流对象转化为安全的Set对象
func NewSafeSetOfStream(s stream.Stream) (ret set.Set) {
	ret = set.NewSafeSetOfLen(s.Size())
	data := s.ToArray()
	ret.Adds(data...)
	return
}

// NewUnSafeSetOfStream 将流对象转化为不安全的Set对象
func NewUnSafeSetOfStream(s stream.Stream) (ret set.Set) {
	ret = set.NewUnsafeSetOfLen(s.Size())
	data := s.ToArray()
	ret.Adds(data...)
	return
}
