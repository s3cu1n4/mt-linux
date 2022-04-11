package main

import "C"
import (
	"mt-linux/common"
	"strings"
)

// 导出函数
//export PrintGo
func PrintGo(s *C.char) {
	datas := strings.Split(C.GoString(s), " ")
	common.DataToMap(datas)
}

//export CheckMod
func CheckMod() {
	common.Rmmod()
}
