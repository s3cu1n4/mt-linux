package common

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/s3cu1n4/logs"
)

var (
	DOWNLOAD_HOSTS = []string{
		"https://lf3-elkeid.bytetos.com/obj/elkeid-download/ko/",
		"https://lf6-elkeid.bytetos.com/obj/elkeid-download/ko/",
		"https://lf9-elkeid.bytetos.com/obj/elkeid-download/ko/",
		"https://lf26-elkeid.bytetos.com/obj/elkeid-download/ko/",
	}
	KMOD_VERSION   = "1.7.0.4"
	KMOD_NAME      = "hids_driver"
	kmodfilename   = "." + KMOD_NAME
	kmod_md5_cache = "." + KMOD_NAME + "md5_cache"
	isinsmod       = false
)

func Checkmod() bool {
	if _, err := os.Stat("/proc/elkeid-endpoint"); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func InstallKO() (err error) {
	checkstatus, err := getKmodcache()
	if err != nil || !checkstatus {
		logs.Error("get kmod cache err:", err)
		// return err
		err = downloadKO()
		if err != nil {
			logs.Error("下载 Elkeid(AgentSmith-HIDS) Driver失败,请参考文档自行编译并加载内核模块: https://github.com/bytedance/Elkeid/blob/main/driver/README-zh_CN.md err:", err)
			return err
		}
	}
	logs.Info("use cached kmod install")
	info, err := insmod(kmodfilename)
	if err != nil {
		logs.Error("insmod err :", info)
		return
	}
	logs.Info("kmod driver install success")
	isinsmod = true
	return
}

func insmod(kmodname string) (output string, err error) {
	out, err := exec.Command("insmod", kmodname).Output()

	if err != nil {
		logs.Fatal(err.Error())
		return
	}
	output = string(out)
	return

}

func Rmmod() (err error) {
	if isinsmod {
		logs.Info("need rmmod hids_driver ")
		// time.Sleep(1 * time.Second)
		command := exec.Command("rmmod", KMOD_NAME)
		err := command.Run()
		if err != nil {
			// log.Fatalf("cmd.run() failed with %s", err)
			logs.Fatalf("rmmod hids_driver err:%s %s", err, string(command.String()))
			return err

		}
		logs.Info("rmmod hids_driver success")

		return err

	} else {
		logs.Info("Unwanted rmmod hids_driver ")
	}
	return

}

func downloadKO() (err error) {
	kname := getkname()
	if kname == "" {
		logs.Fatal("获取内核版本信息失败")
		os.Exit(-1)
	}
	logs.Info("kernel version is: ", kname)
	for _, val := range DOWNLOAD_HOSTS {
		downloadurl := fmt.Sprintf("%s%s_%s_%s_amd64.ko", val, KMOD_NAME, KMOD_VERSION, kname)
		err = downloadkmod(kmodfilename, downloadurl)
		if err != nil {
			logs.Errorf("kmod driver %s 下载失败: %s", downloadurl, err.Error())
			continue
		} else {
			return
		}
	}
	return

}

func downloadkmod(filepath string, url string) (err error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	filemd5 := strings.Trim(resp.Header.Get("Etag"), "\"")

	// Create the file
	out, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)

	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logs.Error("check kmod driver md5 err:", err.Error())
		return

	}

	checkmd5, err := Md5sum(filepath, 1)
	if err != nil {
		logs.Error("check kmod driver md5 err:", err.Error())
		return

	}
	if checkmd5 != filemd5 {
		logs.Info("kmod driver Md5 check is failed")
		return errors.New("md5 check is failed")
	}

	err = WriteToFile(kmod_md5_cache, checkmd5)
	if err != nil {
		logs.Error("write kmod driver md5 cache err", err.Error())
		return
	}
	return
}

//获取系统内核版本
func getkname() string {
	out, err := exec.Command("uname", "-r").Output()

	if err != nil {
		logs.Fatal(err.Error())
		return ""
	}
	return strings.Trim(string(out), "\n")
}

func getKmodcache() (checkstatus bool, err error) {
	contents, err := ioutil.ReadFile(kmod_md5_cache)
	if err != nil {
		return false, err

	} else {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		data := strings.Replace(string(contents), "\n", "", 1)
		checkmd5, err := Md5sum(kmodfilename, 1)
		if err != nil {
			return false, err
		} else if data == checkmd5 {
			return true, nil

		} else {
			return false, errors.New("kmod driver md5 check err")
		}
	}

}
