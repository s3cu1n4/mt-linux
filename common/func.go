package common

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"

	"github.com/s3cu1n4/logs"
)

func SliceToMap(s_key []string, s_value [][]byte, datatype int64) {

	var info string

	for s_key_index := range s_key {
		info = info + fmt.Sprintf("%s:%s, ", s_key[s_key_index], string(s_value[s_key_index]))
	}
	info = strings.Trim(info, ", ")
	logs.Infof("DataType:%s, %s", DataType[datatype], info)
}

func DataToMap(data [][]byte) {
	inttype, _ := strconv.ParseInt(string(data[0]), 10, 64)

	datastruct := DataStruct[inttype]
	if datastruct != nil {
		SliceToMap(datastruct, data[1:], inttype)
	} else {
		logs.Error("NotSliceToMap:", data)
	}
}

//checktype 1 is check file md5
//checktype 2 is check string md5
func Md5sum(filepath string, checktype int) (md5str string, err error) {
	if checktype == 1 {
		f, err := os.Open(filepath)
		if err != nil {
			str1 := "Open err"
			return str1, err
		}
		defer f.Close()

		body, err := ioutil.ReadAll(f)
		if err != nil {
			str2 := "ioutil.ReadAll"
			return str2, err
		}
		md5str = fmt.Sprintf("%x", md5.Sum(body))
		runtime.GC()
	} else if checktype == 2 {
		data := []byte(filepath)
		has := md5.Sum(data)
		md5str = fmt.Sprintf("%x", has)
	}
	return md5str, nil
}

func WriteToFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("file create failed. err: " + err.Error())
	} else {
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt([]byte(content), n)
		fmt.Println("write succeed!")
		defer f.Close()
	}
	return err
}

func CheckEnvironment() bool {
	sysType := runtime.GOOS
	if runtime.GOARCH != "amd64" {
		logs.Error("本程序只支持amd64架构的linux系统使用")
		return false

	}

	if sysType == "linux" {
		user, err := user.Current()
		if err != nil {
			logs.Error("获取系统用户信息失败，error:", err.Error())
			log.Fatalf(err.Error())
		}
		if user.Uid != "0" {
			logs.Errorf("当前用户为:%s,请在root用户下使用本程序", user.Username)
			return false
		}
		return true

	} else {
		logs.Error("本程序只支持Linux系统使用")
		return false
	}

}
