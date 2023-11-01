package common

import (
	"fmt"
	"github.com/mitchellh/go-ps"
	"golang.org/x/exp/slog"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"testing"
	"time"
	"unsafe"
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
	if results == nil {
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

func TestLog(t *testing.T) {
	logFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		t.Fatal(err)
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))
	slog.Info("hello", "count", 3)
}

func TestWalk(t *testing.T) {
	files, err := os.ReadDir("../")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.IsDir())
	}
}

var (
	psapidll                = syscall.NewLazyDLL("psapi.dll")
	procGetModuleFileNameEx = psapidll.NewProc("GetModuleFileNameExW")
)

func getProcessPath(pid uint32) string {
	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, pid)
	if err != nil {
		fmt.Printf("Failed to open process: %s\n", err)
		return ""
	}
	defer syscall.CloseHandle(handle)

	var path [syscall.MAX_PATH]uint16
	ret, _, err := procGetModuleFileNameEx.Call(uintptr(handle), 0, uintptr(unsafe.Pointer(&path[0])), syscall.MAX_PATH)
	if ret == 0 {
		fmt.Println("Failed to get module file name:", err)
		return ""
	}

	absolutePath := syscall.UTF16ToString(path[:])
	fmt.Printf("Process path: %s\n", absolutePath)
	return absolutePath
}

func TestProcess(t *testing.T) {
	cmd := exec.Command("cmd", "/c", "cd D:\\hacker\\tools\\AntSword-Loader-v4.0.3-win32-x64 && AntSword.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		NoInheritHandles: false,
	}
	err := cmd.Start()
	if err != nil {
		fmt.Println("启动出错:", err)
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(3 * time.Second): // 等待一段时间
		// 获取当前系统中所有的进程
		processList, err := ps.Processes()
		if err != nil {
			fmt.Println("获取进程列表失败:", err)
			os.Exit(1)
		}

		// 遍历进程列表，模糊匹配进程名并终止匹配到的进程
		for _, p := range processList {
			if strings.Contains(p.Executable(), "AntSword.exe") {
				absolutePath := getProcessPath(uint32(p.Pid()))
				if strings.HasPrefix(absolutePath, "D:\\hacker\\tools") {
					process, err := os.FindProcess(p.Pid())
					if err != nil {
						return
					}
					err = process.Kill()
					if err != nil {
						fmt.Printf("Failed to terminate process: %s\n", err)
						return
					}
				}
			}
		}
		err = cmd.Process.Kill()
		if err != nil {
			fmt.Println("强制终止进程失败:", err)
		}
		<-done // 等待 cmd.Wait() 完成
		fmt.Println("进程已终止")
	case err := <-done:
		if err != nil {
			fmt.Println("等待进程退出失败:", err)
		} else {
			fmt.Println("进程已退出")
		}
	}
}
