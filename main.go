package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Payload — посылка, которую агент потом будет отправлять на сервер.
type Payload struct {
	Agent   string   `json:"agent"`
	TS      string   `json:"ts"`
	Metrics []Metric `json:"metrics"`
}

func main() {
	host, _ := os.Hostname()

	collectors := []Collector{
		CPUCollector{},
		MemCollector{},
		DiskCollector{},
		NetCollector{},
	}

	payload := Payload{
		Agent: host,
		TS:    time.Now().Format(time.RFC3339),
	}

	for _, c := range collectors {
		metrics, err := c.Collect()
		if err != nil {
			fmt.Fprintf(os.Stderr, "сборщик %s: %v\n", c.Name(), err)
			continue
		}
		for i := range metrics {
			metrics[i].Value = round1(metrics[i].Value)
		}
		payload.Metrics = append(payload.Metrics, metrics...)
	}

	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "JSON:", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
	if err := os.WriteFile("metrics.json", out, 0644); err != nil {
		fmt.Fprintln(os.Stderr, "запись файла:", err)
	}
}
