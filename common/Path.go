package common

import (
	"errors"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
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
	envs := maps.Clone(Paths.Env)
	LoadJava(Paths.Dir)
	LoadPython(Paths.Dir)
	for k, v := range envs {
		Paths.Env[k] = v
	}
	return true, nil
}

func LoadPython(root string) {
	pythons := make([][]string, 0)
	pythonExes := []string{"python.exe", "python3.exe"}
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, exe := range pythonExes {
			if strings.EqualFold(info.Name(), exe) {
				// 找到了 Python 可执行文件，执行命令获取版本信息
				cmd := exec.Command(path, "--version")
				cmd.SysProcAttr = &syscall.SysProcAttr{
					HideWindow: true,
				}
				output, err := cmd.CombinedOutput()
				if err != nil {
					//slog.Error("执行命令出错：%s\n", err)
					return nil
				}

				lines := strings.Split(string(output), "\n")
				version := ""
				if len(lines) > 0 {
					line := lines[0]
					// 假设版本信息是以 "Python" 开头的
					reg, _ := regexp.Compile("Python (\\d+\\.\\d+\\.\\d+)")
					submatch := reg.FindStringSubmatch(line)
					if len(submatch) > 0 {
						version = submatch[1]
					}
				} else {
					version = "未知版本"
				}
				pythons = append(pythons, []string{version, filepath.Dir(path)})
			}
		}

		return nil
	})
	if err != nil {
		slog.Error("遍历目录出错：%s\n", err)
		return
	}
	// 配置python环境变量
	for _, p := range pythons {
		if strings.HasPrefix(p[0], "3.") {
			Paths.Env["python3"] = p[1]
			Paths.Env["python"] = p[1]
			break
		}
	}
}

func LoadJava(root string) {
	javas := make([][]string, 0)
	javaExe := "java.exe"
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.EqualFold(info.Name(), javaExe) {
			// 找到了 java.exe，执行命令获取版本信息
			cmd := exec.Command(path, "-version")
			cmd.SysProcAttr = &syscall.SysProcAttr{
				HideWindow: true,
			}
			output, err := cmd.CombinedOutput()
			if err != nil {
				//fmt.Printf("执行命令出错：%s\n", err)
				return nil
			}
			version := ""
			lines := strings.Split(string(output), "\n")
			if len(lines) > 0 {
				line := lines[0]
				// 假设版本信息是以 "version" 开头的
				reg, _ := regexp.Compile("version \"(.*)\"")
				submatch := reg.FindStringSubmatch(line)
				if len(submatch) > 0 {
					version = submatch[1]
				}
			} else {
				version = "未知版本"
			}
			javas = append(javas, []string{version, filepath.Dir(path)})
		}

		return nil
	})

	if err != nil {
		slog.Error("遍历目录出错：%s\n", err)
		return
	}
	for _, j := range javas {
		//if strings.HasPrefix(j[0], "1.8") {
		//	Paths.Env["java"] = j[1]
		//	continue
		//}
		//处理java8及以下的版本
		r, _ := regexp.Compile("^1\\.(\\d)\\.\\d+")
		ver := r.FindStringSubmatch(j[0])
		if len(ver) == 2 {
			Paths.Env["java"+ver[1]] = j[1]
			continue
		}
		//java8以上自适应
		r, _ = regexp.Compile("(\\d+)\\.\\d+\\.\\d+")
		ver = r.FindStringSubmatch(j[0])
		if len(ver) == 2 {
			Paths.Env["java"+ver[1]] = j[1]
		}

	}
}
