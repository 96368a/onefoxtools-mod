package main

import (
	"changeme/common"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Log("测试开始")
	cmd := &common.Exec{}
	res := cmd.CmdExec("java", "java -jar D:\\BlueTeamTools.jar", "")
	fmt.Println(res)
}
func TestEEE(t *testing.T) {
	cmd := exec.Command("cmd", "/c", "java -jar D:\\BlueTeamTools.jar")
	cmd.Env = append(os.Environ(), "PATH=C:\\Users\\ncxxg\\scoop\\apps\\zulufx11-jdk\\11.66.15\\bin\\")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:       true,
		NoInheritHandles: true,
		CreationFlags:    0x10,
	}

	err := cmd.Start()
	if err != nil {
		fmt.Println("命令启动失败:", err)
		return
	}

	fmt.Println("Process started successfully")

	// 等待进程结束
	//err = cmd.Wait()
	//if err != nil {
	//	fmt.Println("Process finished with error:", err)
	//	os.Exit(1)
	//}

	fmt.Println("Process finished successfully")
}
