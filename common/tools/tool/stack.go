package tool

import (
	"github.com/davecgh/go-spew/spew"
	"runtime"
	"fmt"
)

//PrintPanicStack 产生panic时的调用栈打印
func PrintPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		//logger.Errorf("panic des err:%s ", x.(error).Error())
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			fmt.Println("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			fmt.Println("EXRAS:%d DATA:%s\n", k, spew.Sdump(extras[k]))
		}
	}
}
