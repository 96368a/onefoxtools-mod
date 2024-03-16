package common

import (
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
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
		return err
		//envs, _ := yaml.Marshal(Paths)
		//os.WriteFile("config/base.yml", envs, os.ModePerm)
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

func GenerateConfig() error {
	//CdExePath()
	typeS := `.*wx.StaticBox\(self, wx.ID_ANY, u\"-*([^"]+)-*\"`
	nameS := `self\.(.+?) = wx.Button\(gui.+?u\"(.+?)\"`
	namefile := "GUI_Tools_wxpython_gui.py"
	content, err := os.ReadFile(namefile)
	if err != nil {
		slog.Error("读取文件失败：%s\n", err)
		return err
	}
	bindContent := string(content)
	// 匹配工具类型
	reType := regexp.MustCompile(typeS)
	// 匹配类型下有哪些工具
	reName := regexp.MustCompile(nameS)
	var types []string
	var names [][]string
	var realNames [][]string
	var currentNames []string
	var currentRealNames []string
	for _, line := range strings.Split(bindContent, "\n") {
		if matches := reType.FindStringSubmatch(line); len(matches) > 1 {
			// 匹配类型
			t := strings.ReplaceAll(matches[1], "-", "")
			types = append(types, t)
			if len(currentNames) > 0 {
				names = append(names, currentNames)
				realNames = append(realNames, currentRealNames)
				currentNames = nil
				currentRealNames = nil
			}
		} else if matches := reName.FindStringSubmatch(line); len(matches) > 1 {
			// 匹配工具
			//name := matches[1]
			currentNames = append(currentNames, matches[1])
			currentRealNames = append(currentRealNames, matches[2])
		}
	}
	// 加入最后一个工具
	if len(currentNames) > 0 {
		names = append(names, currentNames)
		realNames = append(realNames, currentRealNames)
	}

	var datas []TypeConfig
	content, err = os.ReadFile("GUI_Tools.py")
	if err != nil {
		slog.Error("读取文件失败：%s\n", err)
		return err
	}
	commandContent := strings.ReplaceAll(string(content), "\n", "")
	i := 0
	for index, n := range names {
		var d []Config
		for ii, name := range n {
			// 匹配绑定事件
			s := `self\.` + name + `\.Bind\(wx.EVT_BUTTON, self\.(.+)\)`
			clickRe := regexp.MustCompile(s)
			clickMatches := clickRe.FindStringSubmatch(bindContent)
			if len(clickMatches) > 1 {
				// 根据绑定事件匹配命令
				clickFunction := clickMatches[1]
				s = `def ` + clickFunction + `\(self, event\):.+?subprocess.Popen\((.+?),.+?\)`
				commandRe := regexp.MustCompile(s)
				commandMatches := commandRe.FindStringSubmatch(commandContent)
				if len(commandMatches) > 1 {
					command := strings.TrimSpace(commandMatches[1])
					command = regexp.MustCompile(`["'+]`).ReplaceAllString(command, " ")
					command = regexp.MustCompile(`\s+`).ReplaceAllString(command, " ")
					command = strings.TrimSpace(command)
					env := ""
					javaRe := regexp.MustCompile(`(java\d{1,2})_path`)
					javaMatches := javaRe.FindStringSubmatch(command)
					if len(javaMatches) > 0 {
						env = javaMatches[1]
						command := javaRe.ReplaceAllString(command, "java")
						d = append(d, Config{
							Name:    realNames[index][ii],
							Command: command,
							Env:     env,
							Index:   ii + 1,
						})
					} else {
						d = append(d, Config{
							Name:    realNames[index][ii],
							Command: command,
							Index:   ii + 1,
						})
					}
					i++
				}
			}
		}
		datas = append(datas, TypeConfig{
			Type:   types[index],
			Config: d,
			Index:  index + 1,
		})
	}

	if err := os.MkdirAll("config/tools", 0777); err != nil {
		slog.Error("创建目录失败:", err)
		return err
	}
	for _, data := range datas {
		filename := fmt.Sprintf("config/tools/%s.yml", data.Type)
		file, err := os.Create(filename)
		if err != nil {
			slog.Error("创建文件失败：%s\n", err)
			return err
		}
		dataBytes, err := yaml.Marshal(data)
		if err != nil {
			slog.Error("转换为YAML失败：%s\n\n", err)
			continue
		}

		if _, err = file.Write(dataBytes); err != nil {
			slog.Error("写入文件失败：%s\n\n", err)
			continue
		}
		slog.Info("写入文件成功：%s\n", filename)
		file.Close()
	}
	LoadEnv(Paths.Dir)
	return nil
}
