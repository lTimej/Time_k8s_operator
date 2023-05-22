package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func getCPU() {
	// cpuInfo, err := cpu.Info()
	// if err != nil {
	// 	fmt.Println("cpu info init err,err: ", err)
	// 	return
	// }
	// for _, cp := range cpuInfo {
	// 	fmt.Printf("cpu: %v,type:%T\n", cp, cp)
	// }
	for {
		n, _ := cpu.Percent(time.Second, false)
		fmt.Println("percent:", n)
	}
}
