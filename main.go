package main

import (
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	host, _ := os.Hostname()
	fmt.Println("Агент:", host)

	// CPU: замер за 1 секунду, общий по всем ядрам
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println("CPU: ошибка:", err)
	} else {
		fmt.Printf("CPU: %.1f %%\n", cpuPercent[0])
	}

	// Оперативная память
	vm, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("RAM: ошибка:", err)
	} else {
		fmt.Printf("RAM: %.1f %% (%d из %d МБ)\n",
			vm.UsedPercent, vm.Used/1024/1024, vm.Total/1024/1024)
	}

	// Диски: агент сам находит все разделы
	parts, err := disk.Partitions(false)
	if err != nil {
		fmt.Println("Диски: ошибка:", err)
		return
	}
	for _, p := range parts {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			fmt.Printf("Диск %s: нет данных (%v)\n", p.Mountpoint, err)
			continue
		}
		fmt.Printf("Диск %s: %.1f %% (свободно %d ГБ)\n",
			p.Mountpoint, usage.UsedPercent, usage.Free/1024/1024/1024)
	}
}
