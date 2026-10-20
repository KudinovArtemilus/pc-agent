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

// collect опрашивает все сборщики и собирает одну посылку.
func collect(name string, collectors []Collector) Payload {
	payload := Payload{
		Agent: name,
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
	return payload
}

func main() {
	path, err := configPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "путь к конфигу:", err)
		os.Exit(1)
	}

	cfg, err := loadConfig(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "конфиг:", err)
		os.Exit(1)
	}

	if cfg.Name == "" {
		cfg.Name, _ = os.Hostname()
	}

	fmt.Fprintf(os.Stderr, "агент %s | интервал %s | сервер %q\n",
		cfg.Name, cfg.Interval, cfg.Server)

	collectors := []Collector{
		CPUCollector{},
		MemCollector{},
		DiskCollector{},
		NetCollector{},
	}

	for {
		payload := collect(cfg.Name, collectors)

		out, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, "JSON:", err)
		} else {
			fmt.Println(string(out))
		}

		time.Sleep(cfg.Interval)
	}
}
