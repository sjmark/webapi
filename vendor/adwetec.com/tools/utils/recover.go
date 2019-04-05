package utils

import "fmt"
// 容错处理
func CatchError(s ...interface{}) {
	fn := func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("CatchError err:%v,s:%v", err, s))
			return
		}
	}
	defer fn()
}