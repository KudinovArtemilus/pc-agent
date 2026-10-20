package main

import (
	"net"
	"time"

	"golang.org/x/sys/windows"
)

// ifCounters — счётчики одного сетевого интерфейса.
// В старом API Windows они 32-битные (uint32).
type ifCounters struct {
	recv, sent    uint32
	errors, drops uint32
}

// readIfCounters снимает счётчики со всех интерфейсов.
// Результат: имя интерфейса -> счётчики.
func readIfCounters() (map[string]ifCounters, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	result := make(map[string]ifCounters)
	for _, ifi := range ifs {
		row := windows.MibIfRow{Index: uint32(ifi.Index)}
		if err := windows.GetIfEntry(&row); err != nil {
			continue // интерфейс не отдал данные — пропускаем только его
		}
		result[ifi.Name] = ifCounters{
			recv:   row.InOctets,
			sent:   row.OutOctets,
			errors: row.InErrors + row.OutErrors,
			drops:  row.InDiscards + row.OutDiscards,
		}
	}
	return result, nil
}

type NetCollector struct{}

func (c NetCollector) Name() string {
	return "network"
}

func (c NetCollector) Collect() ([]Metric, error) {
	before, err := readIfCounters()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second)

	after, err := readIfCounters()
	if err != nil {
		return nil, err
	}

	var metrics []Metric

	for name, cur := range after {
		// служебный интерфейс без трафика — пропускаем
		if cur.recv == 0 && cur.sent == 0 {
			continue
		}

		prev, ok := before[name]
		if !ok {
			continue
		}

		// Разница uint32 правильна даже при переполнении счётчика
		recvBytes := cur.recv - prev.recv
		sentBytes := cur.sent - prev.sent

		labels := map[string]string{"iface": name}

		metrics = append(metrics,
			Metric{Name: "net.recv_kbps", Value: float64(recvBytes) / 1024, Unit: "KB/s", Labels: labels},
			Metric{Name: "net.sent_kbps", Value: float64(sentBytes) / 1024, Unit: "KB/s", Labels: labels},
			Metric{Name: "net.errors", Value: float64(cur.errors), Labels: labels},
			Metric{Name: "net.drops", Value: float64(cur.drops), Labels: labels},
		)
	}

	return metrics, nil
}
