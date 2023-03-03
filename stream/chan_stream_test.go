package stream_test

import (
	"github.com/YiGuan-z/data/stream"
	"testing"
)

var (
	data = []any{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	}
)

func TestNewStream(t *testing.T) {
	s := make([]any, len(data))
	copy(data, s)
	cs := stream.NewChanStream(s)
	as := stream.NewArrayStream(s)
	if cs == nil {
		t.Error("chanStream创建失败")
	}
	if as == nil {
		t.Error("arrayStream创建失败")
	}
}
