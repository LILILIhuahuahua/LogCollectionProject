package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"time"
)

/**
 * @author xhli
 * @date 2021/6/26 20:28
 * @version 1.0
 * @description: TODO
 */

func main(){
	getCpuInfo()
	getCpuLoad()
}

// cpu info
func getCpuInfo() {
	cpuInfos, err := cpu.Info()
	if err != nil {
		fmt.Printf("get cpu info failed, err:%v", err)
	}
	for _, ci := range cpuInfos {
		fmt.Println(ci)
	}
	// CPU使用率
	for {
		percent, _ := cpu.Percent(time.Second, false)
		fmt.Printf("cpu percent:%v\n", percent)
	}
}


func getCpuLoad() {
	info, _ := load.Avg()
	fmt.Printf("%v\n", info)
}