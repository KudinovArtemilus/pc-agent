package main

import "math"

// Metric — одно измерение: имя, значение, единица и метки.
type Metric struct {
	Name   string            `json:"name"`
	Value  float64           `json:"value"`
	Unit   string            `json:"unit,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}

// Collector — любой сборщик метрик.
type Collector interface {
	Name() string
	Collect() ([]Metric, error)
}
