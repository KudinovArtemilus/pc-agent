package main

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUCollector struct{}

func (c CPUCollector) Name() string {
	return "cpu"
}

func (c CPUCollector) Collect() ([]Metric, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}
	return []Metric{
		{Name: "cpu.percent", Value: percent[0], Unit: "%"},
	}, nil
}
