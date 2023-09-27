package common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestCmd(t *testing.T) {
	//newDir := "D:\\hacker\\ONE-FOX集成工具箱_V4.0魔改星球专版_by狐狸"
	t.Log("测试开始")
	exec := &Exec{
		Path: "",
	}
	results := exec.CmdExec("java", "echo %PATH%", "")
	//results := exec.CmdExec("java", "java -jar D:\\BlueTeamTools.jar", "")
	fmt.Println(Paths)
	if results == "" {
		t.Error("测试失败")
	}
	t.Log("测试结束")
}

func TestPowshell(t *testing.T) {

	t.Log("测试开始")
	exec := &Exec{}
	err := exec.PowerShellExec("echo %PATH%")
	if err != nil {
		t.Error("测试失败", err)
	}
	t.Log("测试结束")
}

func TestPath(t *testing.T) {
	// 获取当前工作目录
	currentDir, _ := os.Getwd()
	rootDir := filepath.Dir(currentDir)
	os.Chdir(rootDir)
	t.Log("测试开始")
	fmt.Println(Paths)

}

func getVersionFromOutput(output string) string {
	// 在输出中查找版本信息
	// 假设输出的第一行是包含版本信息的
	lines := strings.Split(output, "\n")
	if len(lines) > 0 {
		line := lines[0]
		// 假设版本信息是以 "version" 开头的
		reg, _ := regexp.Compile("version \"(.*)\"")
		submatch := reg.FindStringSubmatch(line)
		if len(submatch) > 0 {
			return submatch[1]
		}
	}

	return "未知版本"
}
func TestJava(t *testing.T) {
	root := "D:\\sec\\tools\\ONE-FOX集成工具箱_V4.0魔改星球专版_by狐狸" // 当前目录
	javaExe := "java.exe"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.EqualFold(info.Name(), javaExe) {
			// 找到了 java.exe，执行命令获取版本信息
			cmd := exec.Command(path, "-version")
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("执行命令出错：%s\n", err)
				return nil
			}

			version := getVersionFromOutput(string(output))
			fmt.Printf("找到了 java.exe，版本为：%s : %s \n", version, path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("遍历目录出错：%s\n", err)
	}
}
func TestPython(t *testing.T) {
	root := "D:\\sec\\tools\\ONE-FOX集成工具箱_V4.0魔改星球专版_by狐狸" // 当前目录
	pythonExes := []string{"python.exe", "python3.exe"}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
				output, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Printf("执行命令出错：%s\n", err)
					return nil
				}

				//version := getVersionFromOutput(string(output))
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
				}
				fmt.Printf("找到了 %s，版本为：%s : %s\n", exe, version, path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("遍历目录出错：%s\n", err)
	}
}
