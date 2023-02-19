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
