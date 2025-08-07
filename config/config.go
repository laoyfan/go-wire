package config

import (
	"flag"
	"fmt"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name     string
		Mode     string
		Port     int
		Limit    float64
		Burst    int
		ClientID string
	}
	AllowOrigins      []string
	AllowedOriginsMap map[string]struct{}
	Redis             map[string]RedisConfig
	Log               struct {
		Driver     string
		Director   string // 日志文件夹
		Level      string // 日志级别
		MaxAge     int    // 日志保存天数
		MaxSize    int    // 日志大小(MB)
		MaxBackups int    // 日志备份数量
		Format     string // 输出日志格式
	}
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

var ProviderSet = wire.NewSet(NewConfig)

func NewConfig() (*Config, error) {
	// 设置命令行参数
	env := flag.String("env", "debug", "配置环境(debug|prod|test)")
	port := flag.Int("port", 0, "自定义端口（可选）")
	flag.Parse()

	v := viper.New()
	v.AddConfigPath("config") // 指定根目录下的 config 文件夹
	v.SetConfigName(*env)     // 例如 "debug" 对应 debug.yaml
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("配置解析失败: %w", err)
	}

	// 修改端口
	if port != nil && *port > 0 {
		cfg.App.Port = *port
	}

	cfg.AllowedOriginsMap = make(map[string]struct{}, len(cfg.AllowOrigins))
	for _, origin := range cfg.AllowOrigins {
		cfg.AllowedOriginsMap[origin] = struct{}{}
	}

	fmt.Println("配置加载成功", cfg)
	return &cfg, nil
}
