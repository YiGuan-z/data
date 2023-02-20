package set_test

import (
	"fmt"
	"github.com/YiGuan-z/data/set"
	"testing"
)

func TestSet(t *testing.T) {
	t.Run("testAdd", func(t *testing.T) {
		set := set.NewSafeSetOfLen(5)
		set.Add(1)
		set.Add(1)
		set.Add(1)
		set.Add(2)
		set.Range(func(val any) {
			t.Log(val)
		})
		if set.Size() != 2 {
			t.Fatalf("测试失败，元素个数统计不正常，应为2")
		}
	})
	t.Run("testRemove", func(t *testing.T) {
		set := set.NewSafeSetOfLen(5)
		set.Add(1)
		t.Log(set.Size())
		set.Remove(1)
		t.Log(set.Size())
		if set.Size() != 0 {
			t.Fatalf("测试失败，元素应为0个")
		}
	})
	t.Run("testRange", func(t *testing.T) {
		set := set.NewSafeSetOfLen(5)
		type User struct {
			name string
			id   int
		}
		u1 := &User{"小明", 1}
		u2 := &User{"小明", 1}
		u3 := &User{"小明", 1}
		set.Adds(u1, u2, u3)
		set.Add(u1)
		set.Add(u2)
		set.Add(u3)
		t.Log("set中有", set.Size(), "个元素")
		set.Range(func(val any) {
			t.Log(val)
		})
	})
	t.Run("testClone", func(t *testing.T) {
		set := set.NewSafeSetOfLen(5)
		type User struct {
			name string
			id   int
		}
		u1 := &User{"小明", 1}
		u2 := &User{"小明", 1}
		u3 := &User{"小明", 1}
		set.Adds(u1, u2, u3)
		t.Logf("u1%p,u2%p,u3%p\n", u1, u2, u3)
		set.Range(func(val any) {
			t.Logf("%p\r", val)
		})
	})
	t.Run("testClear", func(t *testing.T) {
		set := set.NewSafeSetOfLen(5)
		type User struct {
			name string
			id   int
		}
		u1 := &User{"小明", 1}
		u2 := &User{"小明", 1}
		u3 := &User{"小明", 1}
		set.Adds(u1, u2, u3)
		t.Logf("set size=%d", set.Size())
		set.Clear()
		t.Logf("set size=%d", set.Size())
	})
	t.Run("testCreateNewSafeSetOfSlice", func(t *testing.T) {
		type User struct {
			name string
			id   int
		}
		u1 := User{"小明", 1}
		u2 := User{"小明", 1}
		u3 := User{"小明", 1}
		us := []any{u1, u2, u3}
		s := set.NewSafeSetOfSlice(us)
		t.Logf("set size is %d", s.Size())
		s.Range(func(a any) {
			fmt.Printf("%p\n", &a)
		})
	})
}
