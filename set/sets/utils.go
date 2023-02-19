package sets

import (
	"errors"
	"github.com/YiGuan-z/data/set"
	"github.com/YiGuan-z/data/stream"
)

var (
	newStreamError = errors.New("没有选择需要生成的流")
)

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

func ToArrayStream(set set.Set) (stream.Stream, error) {
	return toStream(set, false, true)
}
func ToChanStream(set set.Set) (stream.Stream, error) {
	return toStream(set, true, false)
}

func NewSetOfStream(s stream.Stream) (ret set.Set) {
	ret = set.NewSet(s.Size())
	s.Range(func(val any) {
		ret.Add(val)
	})
	return
}

func NewSetOfSlice(c []any) set.Set {
	s := set.NewSet(len(c))
	for _, val := range c {
		s.Add(val)
	}
	return s
}
