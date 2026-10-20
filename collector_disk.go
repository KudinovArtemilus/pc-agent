package main

import (
	"github.com/shirou/gopsutil/v3/disk"
)

type DiskCollector struct{}

func (c DiskCollector) Name() string {
	return "disk"
}

func (c DiskCollector) Collect() ([]Metric, error) {
	// 1. Получить список физических разделов
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	// 2. Пустой список, куда будем складывать метрики
	var metrics []Metric

	// 3. Для каждого раздела...
	for _, p := range parts {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue // этот диск пропускаем, идём к следующему
		}

		metrics = append(metrics, Metric{
			Name:   "disk.used_percent",
			Value:  usage.UsedPercent,
			Unit:   "%",
			Labels: map[string]string{"disk": p.Mountpoint},
		})
	}

	// 4. Вернуть всё, что собрали
	return metrics, nil
}
