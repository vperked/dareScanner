package util

import (
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v3/cpu"
)

func MonitorCPU() (float64, error) {
	cpuPercent, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		println("Couldnt fetch cpu:", err.Error())
		return 0, err
	}
	return cpuPercent[0], nil
}

func MonitorMem() (float64, error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		println("Couldnt fetch memory:", err.Error())
		return 0, err
	}
	return mem.UsedPercent, nil
}
