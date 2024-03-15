package common

import (
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
	Title string                `yaml:"title" json:"title"`
	Dir   string                `yaml:"dir" json:"dir"`
	Env   map[string]*EnvConfig `yaml:"env" json:"env"`
}

type EnvConfig struct {
	Current int      `yaml:"current" json:"current"`
	List    []string `yaml:"list" json:"list"`
}

func (i YamlInfo) AppendEnv(name string, env string) {
	if _, ok := i.Env[name]; !ok {
		// 环境变量不存在，则创建
		i.Env[name] = &EnvConfig{
			Current: 0,
			List:    []string{env},
		}

	} else {
		//  检查环境变量是否已经存在
		exists := false
		for _, value := range i.Env[name].List {
			if value == env {
				exists = true
				break
			}
		}
		if !exists {
			i.Env[name].List = append(i.Env[name].List, env)
		}
	}
}

func (i YamlInfo) GetCurrentEnv(name string) string {
	return i.Env[name].List[i.Env[name].Current]
}

func (e EnvConfig) AppendEnv(env string) {
	if e.List == nil {
		e.List = make([]string, 0)
	}
	e.List = append(e.List, env)
}

var Paths YamlInfo

// 初始化环境变量
func InitEnv() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	// 切换到可执行文件所在目录
	os.Chdir(filepath.Dir(executable))
	err = checkConfig()
	if err != nil {
		return err
	}
	data, err := os.ReadFile("config/base.yml")
	if err != nil {
		slog.Error("error:", err)
		return err
	}
	err = yaml.Unmarshal(data, &Paths)
	if err != nil {
		slog.Error("error:", err)
		return err
	}
	// 将配置文件中的相对路径转为绝对路径
	if !filepath.IsAbs(Paths.Dir) {
		Paths.Dir, err = filepath.Abs(Paths.Dir)
		if err != nil {
			slog.Error("error:", err)
			return err
		}
	}
	//envs := maps.Clone(Paths.Env)
	//LoadEnv(Paths.Dir)
	//for k, v := range envs {
	//	Paths.Env[k] = v
	//}
	slog.Info("环境配置加载成功~")
	return nil
}

// 加载java及python环境
func LoadEnv(root string) {
	exes := []string{"java.exe", "python.exe"}
	resultChan := make(chan string)
	maxDepth := 5
	javas := make([][]string, 0)
	pythons := make([][]string, 0)
	// 启动一个goroutine来并发地遍历文件系统
	go func() {
		defer close(resultChan)
		walkDir(root, exes, maxDepth, resultChan)
	}()

	for filePath := range resultChan {
		if filepath.Base(filePath) == "java.exe" {
			java, err := getJavaVersion(filePath)
			if err != nil {
				continue
			}
			javas = append(javas, java)
		} else if filepath.Base(filePath) == "python.exe" {
			python, err := getPythonVersion(filePath)
			if err != nil {
				continue
			}
			pythons = append(pythons, python)
		}
	}
	// 配置java环境变量
	for _, j := range javas {
		//处理java8及以下的版本
		r, _ := regexp.Compile("^1\\.(\\d)\\.\\d+")
		ver := r.FindStringSubmatch(j[0])
		if len(ver) == 2 {
			Paths.AppendEnv("java"+ver[1], j[1])
			continue
		}
		//java8以上自适应
		r, _ = regexp.Compile("(\\d+)\\.\\d+\\.\\d+")
		ver = r.FindStringSubmatch(j[0])
		if len(ver) == 2 {
			Paths.AppendEnv("java"+ver[1], j[1])
		}
	}
	// 配置python环境变量
	for _, p := range pythons {
		if strings.HasPrefix(p[0], "3.") {
			Paths.Env["python3"].AppendEnv(p[1])
			Paths.Env["python"].AppendEnv(p[1])
			break
		}
	}
	marshal, err := yaml.Marshal(Paths)
	if err != nil {
		return
	}
	os.WriteFile("config/base.yml", marshal, 0644)
	slog.Info("环境变量保存成功", Paths.Env)
}

func getPythonVersion(path string) ([]string, error) {
	// 找到了 Python 可执行文件，执行命令获取版本信息
	cmd := exec.Command(path, "--version")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		//slog.Error("执行命令出错：%s\n", err)
		return nil, err
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
	return []string{version, filepath.Dir(path)}, nil
}

func getJavaVersion(path string) ([]string, error) {
	// 执行java.exe获取版本信息
	cmd := exec.Command(path, "-version")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
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
	return []string{version, filepath.Dir(path)}, nil
}

func walkDir(dir string, fileName []string, maxDepth int, resultChan chan<- string) {
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			// 检查子目录级别
			depth := getDepth(dir, path)
			if depth > maxDepth {
				return filepath.SkipDir
			}
		}

		// 检查文件名是否匹配
		for _, f := range fileName {
			if d.Name() == f {
				resultChan <- path
			}
		}
		return nil
	})
}

func getDepth(root string, path string) int {
	depth := 0
	relPath, _ := filepath.Rel(root, path)
	for _, c := range relPath {
		if c == filepath.Separator {
			depth++
		}
	}
	return depth
}

func CdExePath() error {
	// 获取可执行文件路径，读取目录下配置文件
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	dir := filepath.Dir(executable)
	os.Chdir(dir)
	return nil
}
