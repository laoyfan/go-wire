package config

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"path/filepath"
)

type Config struct {
	App struct {
		Name string
		Mode string
		Port int
	}
	Redis map[string]RedisConfig
	Log   struct {
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

//go:embed *.yaml
var configFS embed.FS

var ProviderSet = wire.NewSet(NewConfig)

func NewConfig() (*Config, error) {
	// 设置命令行参数
	env := flag.String("env", "debug", "配置环境(debug|prod|test)")
	port := flag.Int("port", 0, "自定义端口（可选）")
	flag.Parse()
	fileName := fmt.Sprintf("%s.yaml", *env)

	// 读取嵌入的配置文件
	data, err := configFS.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件 %s: %w", fileName, err)
	}

	v := viper.New()
	v.SetConfigType(filepath.Ext(fileName)[1:])

	if err = v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	var cfg Config
	if err = v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("配置解析失败: %w", err)
	}

	// 修改端口
	if port != nil && *port > 0 {
		cfg.App.Port = *port
	}

	fmt.Println("配置加载成功", cfg)
	return &cfg, nil
}
