package main

import (
	"changeme/common"
	"context"
	"errors"
	"github.com/labstack/gommon/log"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CONFIG struct {
	ctx context.Context
}

type Config struct {
	Index   int    `json:"index" yaml:"index"`
	Name    string `json:"name" yaml:"name"`
	Command string `json:"command" yaml:"command"`
	Env     string `json:"env" yaml:"env,omitempty"`
	Dir     string `json:"dir" yaml:"dir,omitempty"`
}

type TypeConfig struct {
	Index  int      `yaml:"index,omitempty" json:"index"`
	Type   string   `yaml:"type" json:"type"`
	Config []Config `yaml:"config" json:"config"'`
}

var Configs []TypeConfig

func (c CONFIG) Start(config Config) bool {
	cmd := &common.Exec{}
	log.Info("执行命令:", config)
	cmd.CmdExec(config.Env, config.Command, config.Dir)
	return true
}

func (c CONFIG) GetConfigs() ([]TypeConfig, error) {
	return Configs, nil
}

func (c CONFIG) GetENVConfigs() (common.YamlInfo, error) {
	return common.Paths, nil
}

func (c CONFIG) GetStartTime() time.Time {
	return startTime
}

func (c CONFIG) GetRefreshTime() time.Time {
	return refreshTime
}

func (c CONFIG) InitConfig() (bool, error) {
	refreshTime = time.Now()
	_, err := InitConfig()
	if err != nil {
		return false, err
	}
	common.InitEnv()
	return true, nil
}

func InitConfig() (bool, error) {
	// 获取可执行文件路径，读取目录下配置文件
	executable, err := os.Executable()
	if err != nil {
		return false, err
	}
	dir := filepath.Dir(executable)
	os.Chdir(dir)
	Configs = make([]TypeConfig, 0)
	// 遍历config下的yaml文件
	err = filepath.WalkDir("config", func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 检查文件扩展名是否为 .yml 或 .yaml
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".yml" || ext == ".yaml" {
			var typeConfig TypeConfig
			data, err := os.ReadFile(path)
			if err != nil {
				slog.Error("error:", err)
			}
			err = yaml.Unmarshal(data, &typeConfig)
			if err != nil {
				slog.Error("error:", err)
				return nil
			}
			Configs = append(Configs, typeConfig)
			slog.Info("配置加载成功:", path)
		}

		return nil
	})

	if err != nil {
		slog.Error("error: 配置文件夹不存在")
		return false, errors.New("配置文件夹不存在")
	}
	os.Chdir(common.Paths.Dir)
	return true, nil
}

func InitLog() {
	//logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	//if err != nil {
	//	slog.Error("error:", err)
	//}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	})))
}
