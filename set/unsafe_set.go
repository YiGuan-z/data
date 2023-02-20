package set

type unsafeSet struct {
	//用于存储数据
	data map[any]struct{}
	//用于统计元素的个数
	size int
}

func (u *unsafeSet) Add(a any) {
	u.data[a] = struct{}{}
	u.size++
}

func (u *unsafeSet) Adds(a ...any) {
	for _, val := range a {
		u.Add(val)
	}
}

func (u *unsafeSet) Remove(a any) (any, bool) {
	if _, ok := u.data[a]; ok {
		delete(u.data, a)
		return u, true
	}
	return nil, false
}

func (u *unsafeSet) Range(f func(any)) {
	for k := range u.data {
		f(k)
	}
}

func (u *unsafeSet) Filter(f func(any) bool) Set {
	u.Range(func(a any) {
		ok := f(a)
		if !ok {
			delete(u.data, a)
			u.size--
		}
		//修正长度
		if u.size != len(u.data) {
			u.size = len(u.data)
		}
	})
	return u
}

func (u *unsafeSet) Find(f func(any) bool) (ret *interface{}) {
	if f == nil {
		panic("没有定义查找方法")
	}
	u.Range(func(val any) {
		ok := f(val)
		if ok {
			ret = &val
			return
		}
	})
	ret = nil
	return
}

func (u *unsafeSet) Clear() {
	if u.Size() > 0 {
		u.Range(func(val any) {
			delete(u.data, val)
			u.size--
		})
	}
}

func (u *unsafeSet) Size() int {
	return u.size
}

func (u *unsafeSet) IsEmpty() bool {
	return 0 == u.size
}

func (u *unsafeSet) Contains(a any) bool {
	_, ok := u.data[a]
	return ok
}

func (u *unsafeSet) ContainsAll(anies []any) bool {
	r := true
	for _, val := range anies {
		_, ok := u.data[val]
		r = r && ok
	}
	return r
}

func (u *unsafeSet) AddAll(anies []any) bool {
	for _, val := range anies {
		u.Add(val)
	}
	return true
}

func (u *unsafeSet) RetainAll(anies []any) {
	for _, val := range anies {
		if value, ok := u.data[val]; !ok {
			u.Remove(value)
		}
	}
}

func (u *unsafeSet) RemoveAll(anies []any) {
	for _, val := range anies {
		if value, ok := u.data[val]; ok {
			u.Remove(value)
		}
	}
}

func (u *unsafeSet) ToSlice() []any {
	ret := make([]any, u.Size())
	u.Range(func(val any) {
		ret = append(ret, val)
	})
	defer u.Clear()
	return ret
}
