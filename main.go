package main

import (
	"fmt"
	"mt-linux/common"
	"os"
	"sync"
	"syscall"

	"github.com/s3cu1n4/logs"
)

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
						fmt.Println(string(data[:i]))
						// dataChan <- data
					}
					s = i + i
				} else if data[i] == 0x1e {
					data[i] = ','
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			data, ok := <-dataChan
			if !ok {
				return
			} else {
				fmt.Println(string(data))
			}

			fmt.Println("asdfasdfasfd")
		}
	}()
	wg.Wait()
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
