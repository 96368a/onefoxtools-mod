package common

import (
	"fmt"
	"os"
	"path/filepath"
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
