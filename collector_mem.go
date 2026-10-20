package main

import (
	"github.com/shirou/gopsutil/v3/mem"
)

type MemCollector struct{}

func (c MemCollector) Name() string {
	return "memory"
}

func (c MemCollector) Collect() ([]Metric, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	return []Metric{
		{Name: "mem.used_percent", Value: vm.UsedPercent, Unit: "%"},
		{Name: "mem.used_mb", Value: float64(vm.Used) / 1024 / 1024, Unit: "MB"},
	}, nil
}
