package main

import (
	"changeme/common"
	"context"
	"github.com/labstack/gommon/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type GOContext struct {
	ctx context.Context
}

var (
	Version   string
	BuildTime string
	Commit    string
)

func (c GOContext) Start(config common.Config) error {
	cmd := &common.Exec{}
	log.Info("执行命令:", config)
	err := cmd.CmdExec(config.Env, config.Command, config.Dir)
	if err != nil {
		return err
	}
	return nil
}

func (c GOContext) TestCmdExec(config common.Config) error {
	cmd := &common.Exec{}
	log.Info("测试命令:", config.Command)
	err := cmd.TestCmdExec(config.Env, config.Command, config.Dir)
	if err != nil {
		return err
	}
	return nil
}

func (c GOContext) GetConfigs() ([]common.TypeConfig, error) {
	return common.Configs, nil
}

func (c GOContext) GetENVConfigs() (common.YamlInfo, error) {
	return common.Paths, nil
}

func (c GOContext) SaveENVConfigs(env common.YamlInfo) error {
	envs, err := yaml.Marshal(env)
	if err != nil {
		return err
	}
	//改变目录
	common.CdExePath()
	//保存配置文件
	os.WriteFile("config/base.yml", envs, os.ModePerm)
	return nil
}

func (c GOContext) GetStartTime() time.Time {
	return startTime
}

func (c GOContext) GetRefreshTime() time.Time {
	return refreshTime
}

func (c GOContext) InitEnv() error {
	return common.InitEnv()
}
func (c GOContext) InitConfig() error {
	refreshTime = time.Now()
	if err := common.InitEnv(); err != nil {
		// 配置文件不存在，跳转到初始化页面
		//if errors.Is(err, os.ErrNotExist) {
		//	runtime.EventsEmit(*wailsContext, "navigate", "/init")
		//} else {
		return err
		//}
	}
	return common.InitConfig()
}
func (c GOContext) GenerateConfig() error {
	err := common.GenerateConfig()
	if err != nil {
		return err
	}
	// 根据配置文件更改应用标题
	if common.Paths.Title != "" {
		runtime.WindowSetTitle(*wailsContext, common.Paths.Title)
	}
	return nil
}
func (c GOContext) Exit() {
	os.Exit(0)
}

func (c GOContext) GetVersion() string {
	return Version
}

func (c GOContext) GetBuildTime() string {
	return BuildTime
}

func (c GOContext) GetGitCommit() string {
	return Commit
}
