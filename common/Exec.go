package common

import (
	"context"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
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

func (e *Exec) CmdExec(env string, command string, workDir string) string {
	// 改变当前工作目录
	_ = os.Chdir(workDir)
	// 开启一个cmd

	cmd := exec.Command("cmd", "/c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:       true,
		NoInheritHandles: true,
		CreationFlags:    0x10,
	}
	if env != "" {
		if _, ok := Paths.Env[env]; ok {
			cmd.Env = append(os.Environ(), "PATH="+Paths.Env[env])
		}
	}
	//for p := range Paths {
	//	if Paths[p].Name == env {
	//		cmd.Env = append(os.Environ(), "PATH="+Paths[p].Path)
	//	}
	//}
	err := cmd.Start()
	if err != nil {
		fmt.Println("命令启动失败:", err)
		return ""
	}
	//cmd.Stdin = strings.NewReader(utf16Command)
	//output, _ := cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Println("执行命令出错:", err)
	//	return err.Error()
	//}
	//results := string(output)
	//fmt.Println("命令输出结果:", results)
	//return results
	//err = cmd.Wait() // 等待命令执行完成
	//if err != nil {
	//	fmt.Println(err)
	//}
	return ""
}
