package main

/*
#include "func.h"
*/
import "C"
import (
	"mt-linux/common"
	"os"

	"github.com/s3cu1n4/logs"
)

func Startsysmem() {
	C.Getsysmem()

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

	Startsysmem()
}
