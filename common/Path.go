package common

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type PathInfo struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type YamlInfo struct {
	Dir string            `yaml:"dir"`
	Env map[string]string `yaml:"env"`
}

var Paths YamlInfo

func InitEnv() {
	data, err := ioutil.ReadFile("env.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &Paths)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	log.Info("init:", Paths)
	if Paths.Dir != "" {
		stat, err := os.Stat(Paths.Dir)
		if err != nil {
			return
		}
		if stat.IsDir() {
			os.Chdir(Paths.Dir)
		} else {
			log.Fatal("目录不存在", Paths.Dir)
		}
	}
	//os.Chdir(Paths.Dir)
	//fmt.Println("init:", Paths)
}
