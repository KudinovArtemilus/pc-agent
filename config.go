package main

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config — настройки агента из agent.yaml.
type Config struct {
	Name     string        `yaml:"name"`
	Server   string        `yaml:"server"`
	Interval time.Duration `yaml:"interval"`
}

// configPath возвращает путь к agent.yaml рядом с agent.exe.
func configPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exe), "agent.yaml"), nil
}

// loadConfig читает конфиг. Нет файла — не ошибка: берём значения по умолчанию.
func loadConfig(path string) (Config, error) {
	cfg := Config{
		Interval: 10 * time.Second,
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
