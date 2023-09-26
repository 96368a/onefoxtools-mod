package common

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
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
	//os.Chdir(Paths.Dir)
	//fmt.Println("init:", Paths)
}
