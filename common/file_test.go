package common

import (
	"fmt"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"testing"
)

func BenchmarkWalk0(b *testing.B) {
	root := "D:\\hacker\\ONE-FOX集成工具箱_V4.0魔改星球专版_by狐狸"
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
			fmt.Println(version)
		}

		return nil
	})

	if err != nil {
		slog.Error("遍历目录出错：%s\n", err)
		return
	}
	//for _, j := range javas {
	//	//if strings.HasPrefix(j[0], "1.8") {
	//	//	Paths.Env["java"] = j[1]
	//	//	continue
	//	//}
	//	//处理java8及以下的版本
	//	r, _ := regexp.Compile("^1\\.(\\d)\\.\\d+")
	//	ver := r.FindStringSubmatch(j[0])
	//	if len(ver) == 2 {
	//		Paths.Env["java"+ver[1]] = j[1]
	//		continue
	//	}
	//	//java8以上自适应
	//	r, _ = regexp.Compile("(\\d+)\\.\\d+\\.\\d+")
	//	ver = r.FindStringSubmatch(j[0])
	//	if len(ver) == 2 {
	//		Paths.Env["java"+ver[1]] = j[1]
	//	}
	//
	//}
}

func BenchmarkWalk1(b *testing.B) {
	rootDir := "D:\\hacker\\ONE-FOX集成工具箱_V4.0魔改星球专版_by狐狸" // 根目录
	fileName := "java.exe"                                // 目标文件名
	maxDepth := 5                                         // 最大子目录级别
	javas := make([][]string, 0)
	// 创建一个channel用于接收匹配到的文件路径
	resultChan := make(chan string)

	// 启动一个goroutine来并发地遍历文件系统
	go func() {
		defer close(resultChan)
		walkDir1(rootDir, fileName, maxDepth, resultChan)
	}()

	// 从channel中读取匹配到的文件路径并打印
	for filePath := range resultChan {
		cmd := exec.Command(filePath, "-version")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			//fmt.Printf("执行命令出错：%s\n", err)
			continue
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
		javas = append(javas, []string{version, filepath.Dir(filePath)})
		fmt.Println(version)
	}
}

func walkDir1(dir, fileName string, maxDepth int, resultChan chan<- string) {
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// 处理错误
			fmt.Printf("Error accessing path %s: %v\n", path, err)
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
		if d.Name() == fileName {
			resultChan <- path
		}

		return nil
	})
}

func TestParse(t *testing.T) {
	os.Chdir("D:\\code_field\\go_code\\wails\\OneFoxTools\\build\\bin\\test")
	types := []string{}
	names := [][]string{}
	typeS := `.*wx.StaticBox\(self, wx.ID_ANY, u\"-*([^"]+)-*\"`
	nameS := `self\.(.+?) = wx.Button\(gui.+?u\"(.+?)\"`
	namefile := ""
	if _, err := os.Stat("GUI_Tools_name.py"); err == nil {
		namefile = "GUI_Tools_name.py"
	} else if _, err := os.Stat("GUI_Tools_wxpython_gui.py"); err == nil {
		namefile = "GUI_Tools_wxpython_gui.py"
	} else {
		fmt.Println("没有找到GUI_Tools_name.py或GUI_Tools_wxpython_gui.py")
		os.Exit(0)
	}
	content, err := ioutil.ReadFile(namefile)
	if err != nil {
		fmt.Println("读取文件失败:", err)
		os.Exit(1)
	}
	c := string(content)
	lines := strings.Split(string(content), "\n")
	n := []string{}
	for _, line := range lines {
		if matched, _ := regexp.MatchString(typeS, line); matched {
			re := regexp.MustCompile(typeS)
			t := re.FindStringSubmatch(line)[1]
			types = append(types, t)
			if len(n) > 0 {
				names = append(names, n)
			}
			n = []string{}
		}
		if matched, _ := regexp.MatchString(nameS, line); matched {
			re := regexp.MustCompile(nameS)
			name := re.FindStringSubmatch(line)
			n = append(n, name[1])
		}
	}
	names = append(names, n)

	datas := []TypeConfig{}
	content, err = ioutil.ReadFile("GUI_Tools.py")
	if err != nil {
		fmt.Println("读取文件失败:", err)
		os.Exit(1)
	}
	cc := string(content)
	cc = strings.ReplaceAll(cc, "\n", "")
	i := 0
	for index, n := range names {
		d := []Config{}
		for _, name := range n {
			s := `self\.` + name + `\.Bind\(wx.EVT_BUTTON, self\.(.+)\)`
			re := regexp.MustCompile(s)
			click := re.FindStringSubmatch(c)
			s = `def ` + click[1] + `\(self, event\):.+?subprocess.Popen\((.+?),.+?\)`
			re = regexp.MustCompile(s)
			rr := re.FindStringSubmatch(cc)
			if len(rr) > 0 {
				command := rr[1]
				command = strings.TrimSpace(command)
				command = strings.ReplaceAll(command, `\"`, "")
				command = strings.ReplaceAll(command, `'`, "")
				command = regexp.MustCompile(`\s+`).ReplaceAllString(command, " ")
				env := ""
				java := regexp.MustCompile(`(java\d{1,2})_path`).FindStringSubmatch(command)
				if len(java) > 0 {
					env = java[1]
					command = strings.ReplaceAll(command, `java\d{1,2}_path`, "java")
					d = append(d, Config{
						Name:    name,
						Command: command,
						Env:     env,
					})
				} else {
					d = append(d, Config{
						Name:    name,
						Command: command,
					})
				}
				fmt.Println(command)
				i++
			}
		}
		datas = append(datas, TypeConfig{
			Type:   types[index],
			Config: d,
		})
	}
	fmt.Println(i)

	if _, err := os.Stat("config"); os.IsNotExist(err) {
		os.Mkdir("config", 0755)
	}
	for _, data := range datas {
		filename := fmt.Sprintf("config/%s.yml", data.Type)
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("创建文件失败:", err)
			os.Exit(1)
		}
		defer file.Close()

		dataBytes, err := yaml.Marshal(data)
		if err != nil {
			fmt.Println("转换为YAML失败:", err)
			os.Exit(1)
		}

		_, err = file.Write(dataBytes)
		if err != nil {
			fmt.Println("写入文件失败:", err)
			os.Exit(1)
		}
	}
	fmt.Println("完成")
}

func TestParse1(t *testing.T) {
	os.Chdir("D:\\code_field\\go_code\\wails\\OneFoxTools\\build\\bin\\test")

	fileName, err := getFileName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("读取文件失败:", err)
		os.Exit(1)
	}

	types, names, err := parseContent(string(content))
	if err != nil {
		fmt.Println("解析内容失败:", err)
		os.Exit(1)
	}

	datas, err := processData(types, names, string(content))
	if err != nil {
		fmt.Println("处理数据失败:", err)
		os.Exit(1)
	}

	err = saveData(datas)
	if err != nil {
		fmt.Println("保存数据失败:", err)
		os.Exit(1)
	}

	fmt.Println("完成")
}

func getFileName() (string, error) {
	fileNames := []string{"GUI_Tools_name.py", "GUI_Tools_wxpython_gui.py"}
	for _, fileName := range fileNames {
		if _, err := os.Stat(fileName); err == nil {
			return fileName, nil
		}
	}
	return "", fmt.Errorf("没有找到GUI_Tools_name.py或GUI_Tools_wxpython_gui.py")
}

func parseContent(content string) ([]string, [][]string, error) {
	var types []string
	var names [][]string
	var currentNames []string

	typeS := regexp.MustCompile(`.*wx.StaticBox\(self, wx.ID_ANY, u\"-*([^"]+)-*\"`)
	nameS := regexp.MustCompile(`self\.(.+?) = wx.Button\(gui.+?u\"(.+?)\"`)

	for _, line := range strings.Split(content, "\n") {
		if typeS.MatchString(line) {
			types = append(types, typeS.FindStringSubmatch(line)[1])
			if len(currentNames) > 0 {
				names = append(names, currentNames)
			}
			currentNames = []string{}
		}
		if nameS.MatchString(line) {
			name := nameS.FindStringSubmatch(line)
			currentNames = append(currentNames, name[1])
		}
	}
	names = append(names, currentNames)
	return types, names, nil
}

func processData(types []string, names [][]string, content string) ([]TypeConfig, error) {
	var datas []TypeConfig

	for index, currentNames := range names {
		var configs []Config
		for _, name := range currentNames {
			config, err := processName(name, content)
			if err != nil {
				return nil, err
			}
			configs = append(configs, config)
		}
		datas = append(datas, TypeConfig{
			Type:   types[index],
			Config: configs,
		})
	}
	return datas, nil
}

func processName(name, content string) (Config, error) {
	bindS := regexp.MustCompile(`self\.` + name + `\.Bind\(wx.EVT_BUTTON, self\.(.+)\)`)
	commandS := regexp.MustCompile(`def ` + bindS.FindStringSubmatch(content)[1] + `\(self, event\):.+?subprocess.Popen\((.+?),.+?\)`)
	javaS := regexp.MustCompile(`(java\d{1,2})_path`)

	command := commandS.FindStringSubmatch(content)[1]
	command = strings.TrimSpace(command)
	command = strings.ReplaceAll(command, `\"`, "")
	command = strings.ReplaceAll(command, `'`, "")
	command = regexp.MustCompile(`\s+`).ReplaceAllString(command, " ")

	var env string
	java := javaS.FindStringSubmatch(command)
	if len(java) > 0 {
		env = java[1]
		command = strings.ReplaceAll(command, `java\d{1,2}_path`, "java")
	}

	return Config{
		Name:    name,
		Command: command,
		Env:     env,
	}, nil
}

func saveData(datas []TypeConfig) error {
	err := os.MkdirAll("config", 0755)
	if err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	for _, data := range datas {
		err := saveDataToFile(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveDataToFile(data TypeConfig) error {
	fileName := fmt.Sprintf("config/%s.yml", data.Type)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	dataBytes, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("转换为YAML失败: %w", err)
	}

	_, err = file.Write(dataBytes)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}
