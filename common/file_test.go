package common

import (
	"fmt"
	"golang.org/x/exp/slog"
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
