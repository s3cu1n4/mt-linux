package main

import (
	"bytes"
	"mt-linux/common"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/s3cu1n4/logs"
)

var exitinfo bool

func readhookdata() {
	fd, err := syscall.Open("/proc/elkeid-endpoint", syscall.O_RDONLY, 0)
	if err != nil {
		logs.Error("Failed on open elkeid file: ", err)
	}
	defer syscall.Close(fd)

	var wg sync.WaitGroup

	wg.Add(2)

	dataChan := make(chan []byte, 100)
	go func() {
		wg.Done()
		for {
			data := make([]byte, 65536)
			s := 0
			n, err := syscall.Read(fd, data)
			if err != nil {
				logs.Error("read fd error", err.Error())
			}
			for i := 0; i <= n; i++ {
				if i == n || data[i] == 0x17 {
					data[i] = 0
					if i > s {
						dataChan <- data[:i]
					}
					s = i + i
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			data, ok := <-dataChan
			if !ok || exitinfo {
				wg.Done()
				return
			} else {
				b := bytes.Split(data, []byte{0x1e})
				common.DataToMap(b)
			}
		}
	}()
	go exitsignal(fd)
	wg.Wait()
}

func exitsignal(fd int) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logs.Info("get signal: ", sig)
		done <- true
	}()
	logs.Info("Server Start Awaiting Signal")

	<-done
	//close fd
	err := syscall.Close(fd)
	if err != nil {
		logs.Error(err)

	}
	exitinfo = true
	common.Rmmod()

	logs.Info("Exit sucess")
}

func main() {
	syscheck := common.CheckEnvironment()
	if !syscheck {
		os.Exit(-1)
	}
	if !common.Checkmod() {
		logs.Info("AgentSmith-HIDS Driver is not installed")
		err := common.InstallKO()
		if err != nil {
			logs.Error("install AgentSmith-HIDS Driver err: ", err)
			os.Exit(-1)
		}
	}
	readhookdata()

}
