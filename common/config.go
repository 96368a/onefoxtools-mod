package common

import (
	"errors"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

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
	Config []Config `yaml:"config" json:"config"`
}

var Configs []TypeConfig

func checkConfig() error {
	stat, err := os.Stat("config/tools")
	if err != nil {
		os.MkdirAll("config/tools", os.ModePerm)
	} else if !stat.IsDir() {
		return errors.New("配置文件夹被占用")
	}

	stat, err = os.Stat("config/base.yml")
	if err != nil {
		envs, _ := yaml.Marshal(Paths)
		os.WriteFile("config/base.yml", envs, os.ModePerm)
	} else if stat.IsDir() {
		return errors.New("基础配置文件被占用")
	}
	return nil
}

func InitConfig() error {
	// 获取可执行文件路径，读取目录下配置文件
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	dir := filepath.Dir(executable)
	os.Chdir(dir)

	// 检查config文件夹
	err = checkConfig()
	if err != nil {
		return err
	}
	Configs = make([]TypeConfig, 0)
	// 遍历config下的yaml文件
	err = filepath.WalkDir("config/tools", func(path string, info os.DirEntry, err error) error {
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
		return errors.New("配置文件夹不存在")
	}
	os.Chdir(Paths.Dir)
	return nil
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
