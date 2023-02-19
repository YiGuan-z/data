package stream_test

import (
	"fmt"
	"github.com/YiGuan-z/data/stream"
	"testing"
)

func TestNewChanStream(t *testing.T) {
	s := []any{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	}
	//stream := NewChanStream(s)
	//stream.Head(2).Range(func(val any) {
	//	fmt.Println(val)
	//})
	stream.NewChanStream(s).Filter(func(val any) bool {
		if i, ok := val.(int); ok {
			return i%2 == 0
		}
		return false
	}).Range(func(val any) {
		fmt.Println(val)
	})
}
