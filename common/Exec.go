package common

import (
	"context"
	"fmt"
	"github.com/mitchellh/go-ps"
	"golang.org/x/exp/slog"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type Exec struct {
	ctx  context.Context
	Path string
}

func GBKtoUTF8(str string) string {
	s, err := simplifiedchinese.GBK.NewDecoder().String(str)
	if err != nil {
		return str
	}
	return s
}

func (e *Exec) PowerShellExec(command string) error {
	// 改变当前工作目录
	_ = os.Chdir(e.Path)
	cmd := exec.Command("powershell.exe")
	cmd.Stdin = strings.NewReader(command + "\n")
	output, err := cmd.Output()
	errStr := cmd.Stderr
	if err != nil {
		fmt.Println("执行命令出错:", err)
		return err
	}
	fmt.Println("命令输出结果:")
	// 将结果写入文件
	ioutil.WriteFile("output.txt", output, 0644)
	fmt.Println(GBKtoUTF8(string(output)))
	fmt.Println(errStr)
	return nil
}

func (e *Exec) CmdExec(env string, command string, workDir string) error {
	// 改变当前工作目录
	_ = os.Chdir(workDir)

	// 开启一个cmd
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second) // 设置超时时间
	defer cancel()

	cmd := exec.Command("cmd", "/c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:       true,
		NoInheritHandles: true,
		CreationFlags:    0x10,
	}
	// 设置环境变量
	if env != "" {
		if _, ok := Paths.Env[env]; ok {
			cmd.Env = append(os.Environ(), "PATH="+Paths.GetCurrentEnv(env))
		}
	}
	err := cmd.Start()
	if err != nil {
		slog.Error("命令启动失败:", err)
		return err
	}
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// 超时后直接返回
		return nil
	case err := <-done:
		if err != nil {
			slog.Error("进程执行出错：", err)
			return err
		} else {
			return nil
		}
	}
}

var (
	psapidll                = syscall.NewLazyDLL("psapi.dll")
	procGetModuleFileNameEx = psapidll.NewProc("GetModuleFileNameExW")
)

func getProcessPath(pid uint32) string {
	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, pid)
	if err != nil {
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
	//fmt.Printf("Process path: %s\n", absolutePath)
	return absolutePath
}
func killProcess(name string, path string) {
	isFind := false
	// 获取当前系统中所有的进程
	processList, err := ps.Processes()
	if err != nil {
		slog.Error("获取进程列表失败:", err)
		return
	}

	// 遍历进程列表，模糊匹配进程名并终止匹配到的进程
	for _, p := range processList {
		if strings.Contains(strings.ToLower(p.Executable()), strings.ToLower(name)) {
			isFind = true
			absolutePath := getProcessPath(uint32(p.Pid()))
			if strings.HasPrefix(absolutePath, path) {
				process, err := os.FindProcess(p.Pid())
				if err != nil {
					return
				}
				err = process.Kill()
				if err != nil {
					slog.Error("Failed to terminate process: %s\n", err)
					return
				} else {
					slog.Error("结束进程成功：", name, path)
				}
			}
		}
	}
	if !isFind {
		slog.Error("未找到进程：", name, path)
	}
}

func (e *Exec) TestCmdExec(env string, command string, workDir string) error {
	// 改变当前工作目录
	_ = os.Chdir(workDir)

	// 开启一个cmd
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second) // 设置超时时间
	defer cancel()

	cmd := exec.Command("cmd", "/c", command)
	// 设置环境变量
	if env != "" {
		if _, ok := Paths.Env[env]; ok {
			cmd.Env = append(os.Environ(), "PATH="+Paths.GetCurrentEnv(env))
		}
	}
	err := cmd.Start()
	if err != nil {
		slog.Error("命令启动失败:", err)
		return err
	}
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// 超时后直接返回
		if env != "" {
			if strings.HasPrefix(env, "java") {
				killProcess("java.exe", Paths.GetCurrentEnv(env))
			} else if strings.HasPrefix(env, "python") {
				killProcess("python.exe", Paths.GetCurrentEnv(env))
			}
		}
		compile := regexp.MustCompile("[\\S]+\\.exe")
		match := compile.FindString(command)
		slog.Info("匹配：", match)
		killProcess(match, Paths.GetCurrentEnv(env))
		cmd.Process.Kill()
		//if err != nil {
		//	slog.Error("进程结束出错：", err)
		//	return err
		//}
		return nil
	case err := <-done:
		if err != nil {
			slog.Error("进程执行出错：", err)
			return err
		} else {
			return nil
		}
	}
}
