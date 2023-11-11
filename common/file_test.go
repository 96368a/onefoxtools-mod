package common

import (
	"fmt"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	"log"
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

	typeS := `.*wx.StaticBox\(self, wx.ID_ANY, u\"-*([^"]+)-*\"`
	nameS := `self\.(.+?) = wx.Button\(gui.+?u\"(.+?)\"`
	namefile := "GUI_Tools_wxpython_gui.py"
	content, err := os.ReadFile(namefile)
	if err != nil {
		log.Fatalf("读取文件失败：%s\n", err)
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
			t := matches[1]
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
		log.Fatalf("读取文件失败：%s\n", err)
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
					javaMatches := regexp.MustCompile(`(java\d{1,2})_path`).FindStringSubmatch(command)
					if len(javaMatches) > 0 {
						env = javaMatches[1]
						command = strings.ReplaceAll(command, `java\d{1,2}_path`, "java")
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
					fmt.Println(command)
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

	//if _, err := os.Stat("config/tools"); os.IsNotExist(err) {
	//	os.MkdirAll("config/tools", 0755)
	//}
	if err := os.MkdirAll("config/tools", 0755); err != nil {
		fmt.Println("创建目录失败:", err)
		os.Exit(1)
	}
	for _, data := range datas {
		filename := fmt.Sprintf("config/tools/%s.yml", data.Type)
		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("创建文件失败：%s\n", err)
		}
		dataBytes, err := yaml.Marshal(data)
		if err != nil {
			log.Printf("转换为YAML失败：%s\n\n", err)
			continue
		}

		if _, err = file.Write(dataBytes); err != nil {
			log.Printf("写入文件失败：%s\n\n", err)
			continue
		}
		file.Close()
	}
}

func TestParse1(t *testing.T) {
	os.Chdir("D:\\code_field\\go_code\\wails\\OneFoxTools\\build\\bin\\test")
	err := GenerateConfig()
	if err != nil {
		return
	}
}
