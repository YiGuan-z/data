package set

/*-------------------------------------safe-------------------------------------*/

// NewSafeSet 创建一个安全Set
func NewSafeSet() Set {
	return &safeSet{
		data: make(map[any]struct{}),
	}
}

// NewSafeSetOfLen 根据length来设置Set的初始化大小
func NewSafeSetOfLen(length int) Set {
	return &safeSet{data: make(map[any]struct{}, length)}
}

// NewSafeSetOfSlice 通过一个切片创建一个安全的Set
func NewSafeSetOfSlice(a []any) Set {
	ret := NewSafeSetOfLen(len(a))
	ret.Adds(a...)
	return ret
}

/*--------------------------------------unsafe--------------------------------------*/

// NewUnsafeSet 创建一个不安全的Set
func NewUnsafeSet() Set {
	return &unsafeSet{
		data: make(map[any]struct{}),
	}
}

// NewUnsafeSetOfLen 通过length来设置Set的初始化大小
func NewUnsafeSetOfLen(length int) Set {
	return &unsafeSet{data: make(map[any]struct{}, length)}
}

// NewUnsafeSetOfSlice 通过一个切片创建一个不安全的Set
func NewUnsafeSetOfSlice(a []any) Set {
	ret := NewUnsafeSetOfLen(len(a))
	ret.Adds(a...)
	return ret
}
