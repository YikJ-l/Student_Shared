package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Config 统一配置结构
// 使用简单的 KEY=VALUE 格式，无分区；键名统一使用大写下划线。
// 示例见 app/conf/congfig.conf
// 支持环境变量覆盖：若存在同名环境变量，优先使用环境变量。

type DBConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

type AIConfig struct {
	APIKey   string
	BaseURL  string
	Model    string
	TimeoutS int
}

type ServerConfig struct {
	Port string
}

type Config struct {
	DB     DBConfig
	AI     AIConfig
	Server ServerConfig
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

	// 先从文件读取 KV
	pairs := map[string]string{}
	if f, err := os.Open(cfgPath); err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") { // 注释/空行
				continue
			}
			// key=value，仅按第一个 '=' 分割
			if i := strings.Index(line, "="); i > 0 {
				k := strings.TrimSpace(line[:i])
				v := strings.TrimSpace(line[i+1:])
				pairs[strings.ToUpper(k)] = v
			}
		}
		_ = f.Close()
	}

	// 环境变量覆盖
	get := func(key string) string {
		if val := os.Getenv(key); val != "" {
			return val
		}
		return pairs[strings.ToUpper(key)]
	}

	// 组装配置，提供合理默认值
	cfg = Config{
		DB: DBConfig{
			User: firstNonEmpty(get("MYSQL_USER"), "root"),
			Pass: firstNonEmpty(get("MYSQL_PASS"), "abc123456"),
			Host: firstNonEmpty(get("MYSQL_HOST"), "localhost"),
			Port: firstNonEmpty(get("MYSQL_PORT"), "3306"),
			Name: firstNonEmpty(get("MYSQL_DB"), "shared_student"),
		},
		AI: AIConfig{
			APIKey:  get("OPENAI_API_KEY"),
			BaseURL: firstNonEmpty(get("OPENAI_BASE_URL"), "https://api.openai.com"),
			Model:   firstNonEmpty(get("OPENAI_MODEL"), "gpt-4o-mini"),
			TimeoutS: atoiSafe(firstNonEmpty(get("OPENAI_TIMEOUT_SECONDS"), "30")),
		},
		Server: ServerConfig{
			Port: firstNonEmpty(get("SERVER_PORT"), "8080"),
		},
	}

	loaded = true
	return cfg
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func atoiSafe(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0
		}
	}
	if s == "" { return 0 }
	// 简单转换，避免引入 strconv 依赖也可使用 strconv.Atoi
	res := 0
	for i := 0; i < len(s); i++ {
		res = res*10 + int(s[i]-'0')
	}
	return res
}