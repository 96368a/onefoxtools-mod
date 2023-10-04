package common

import (
	"errors"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	"os"
)

type PathInfo struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type YamlInfo struct {
	Dir string            `yaml:"dir" json:"dir"`
	Env map[string]string `yaml:"env" json:"env"`
}

var Paths YamlInfo

func InitEnv() (bool, error) {
	data, err := os.ReadFile("env.yml")
	if err != nil {
		slog.Error("error:", err)
		return false, err
	}
	err = yaml.Unmarshal(data, &Paths)
	if err != nil {
		slog.Error("error:", err)
		return false, err
	}
	slog.Info("环境配置加载成功~")
	if Paths.Dir != "" {
		stat, err := os.Stat(Paths.Dir)
		if err != nil {
			return false, errors.New("主目录不存在")
		}
		if stat.IsDir() {
			os.Chdir(Paths.Dir)
		} else {
			return false, errors.New("主目录不是目录")
		}
	}
	return true, nil
}
