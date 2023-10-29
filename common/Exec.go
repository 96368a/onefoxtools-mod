package common

import (
	"context"
	"fmt"
	"golang.org/x/exp/slog"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
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
			cmd.Env = append(os.Environ(), "PATH="+Paths.Env[env])
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
