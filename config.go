package main

import (
	"changeme/common"
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CONFIG struct {
	ctx context.Context
}

type Config struct {
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

var Configs = make([]TypeConfig, 0)

func (c CONFIG) Start(config Config) bool {
	cmd := &common.Exec{}
	log.Info("执行命令:", config)
	cmd.CmdExec(config.Env, config.Command, config.Dir)
	return true
}

func (c CONFIG) GetConfigs() []TypeConfig {
	return Configs
}

func InitConfig() {
	//cmd := exec.Command("cmd", "/C", "cd")
	//result, err := cmd.Output()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//// output res.txt
	//ioutil.WriteFile("res.txt", result, 0644)
	// 遍历config下的yaml文件
	err := filepath.Walk("config", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 检查文件扩展名是否为 .yml 或 .yaml
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".yml" || ext == ".yaml" {
			var typeConfig TypeConfig
			data, err := ioutil.ReadFile(path)
			ioutil.WriteFile("res.txt", []byte(common.Paths.Dir), 0644)
			if err != nil {
				ioutil.WriteFile("res.txt", []byte(err.Error()), 0644)
				fmt.Println("读取文件出错:", err)
				return nil
			}
			err = yaml.Unmarshal(data, &typeConfig)
			if err != nil {
				fmt.Println("error:", err)
				return nil
			}
			Configs = append(Configs, typeConfig)
			//fmt.Println("init:", typeConfig)
			//if _, ok := Configs[typeConfig.Type]; !ok {
			//	Configs[typeConfig.Type] = make([]Config, 0)
			//}
			//Configs[typeConfig.Type] = append(Configs[typeConfig.Type], typeConfig.Config...)
			//
			//fmt.Println("init:", Configs)
		}

		return nil
	})

	if err != nil {
		fmt.Println("遍历文件夹出错:", err)
	}
	//data, err := ioutil.ReadFile("env.yml")
	//if err != nil {
	//	return
	//}
	//err = yaml.Unmarshal(data, &Paths)
	//if err != nil {
	//	fmt.Println("error:", err)
	//	return
	//}
	//fmt.Println("init:", Paths)
	os.Chdir(common.Paths.Dir)
}
