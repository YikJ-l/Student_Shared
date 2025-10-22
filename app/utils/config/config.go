package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type DBConfig struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Port string `json:"port"`
	Name string `json:"name"`
}

type AIConfig struct {
	APIKey   string `json:"apiKey"`
	BaseURL  string `json:"baseUrl"`
	Model    string `json:"model"`
	TimeoutS int    `json:"timeoutSeconds"`
}

type ServerConfig struct {
	Port string `json:"port"`
}

type Config struct {
	DB     DBConfig     `json:"mysql"`
	AI     AIConfig     `json:"openai"`
	Server ServerConfig `json:"server"`
}

var (
	cfg     Config
	loaded  bool
	mu      sync.RWMutex
	cfgPath = filepath.Join("app", "conf", "congfig.conf")
)

// Load 读取配置文件（幂等）。
func Load() Config {
	mu.Lock()
	defer mu.Unlock()
	if loaded {
		return cfg
	}

	if fc, err := os.ReadFile(cfgPath); err == nil {
		_ = json.Unmarshal(fc, &cfg)
	}

	loaded = true
	return cfg
}
