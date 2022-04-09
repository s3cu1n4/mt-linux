package main

import "C"
import (
	"mt-linux/common"
	"strings"
)

// 导出函数
//export PrintGo
func PrintGo(start, end C.int, s *C.char) {
	datas := strings.Split(C.GoString(s), " ")
	// inttype, _ := strconv.ParseInt(datas[0], 10, 64)
	// datatype := common.DataType[inttype]
	// fmt.Printf("%d type: %s data: %s\n", len(datas), datatype, C.GoString(s))
	common.DataToMap(datas)
	// fmt.Printf("start: %d end: %d data: %s\n", start, end, C.GoString(s))
}

//export CheckMod
func CheckMod() {
	common.Rmmod()
}
