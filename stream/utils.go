package stream

// 检查元素是否为空
func checkFunc(f any) {
	if f == nil {
		panic("方法未定义")
	}
}

// newChanel 根据切片创建一个管道并返回大小
func newChanel(data []any) (send <-chan any, lenght int) {
	ret := make(chan any)
	go func() {
		for _, v := range data {
			ret <- v
		}
		close(ret)
	}()
	send = ret
	lenght = len(data)
	return
}
