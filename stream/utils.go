package stream

// checkFunc 检查元素是否为空
func checkFunc(f any) {
	if f == nil {
		panic("方法未定义")
	}
}

func checkSlice(a ...any) bool {
	for _, item := range a {
		if item == nil {
			return false
		}
	}
	return true
}

// newChanel 根据切片创建一个管道并返回大小
func newChanel(data []any) (send <-chan any, lenght int) {
	lenght = len(data)
	if lenght == 0 {
		panic("slice is empty")
	}
	ret := make(chan any)
	go func() {
		for _, v := range data {
			ret <- v
		}
		close(ret)
	}()
	send = ret
	return
}
