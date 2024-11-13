package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type ProgramData struct {
	Name  string
	APPID string
}

var apps []ProgramData

func init() {
	apps, err := getStartApps()

	if err != nil {
		print("err code : 1200")
	}

	print(apps)
}

// 通过 PowerShell 获取所有开始菜单的应用
func getStartApps() ([]ProgramData, error) {
	// PowerShell 命令获取开始菜单应用
	cmd := exec.Command("powershell", "chcp 65001; Get-StartApps | Select Name")

	var out bytes.Buffer
	cmd.Stdout = &out

	// 执行命令并捕获输出
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// 解析 PowerShell 输出
	lines := strings.Split(out.String(), "\n")

	var programs []ProgramData

	// 跳过第一行标题并遍历每一行数据
	for _, line := range lines[4:] { // 跳过标题行

		// print(line + "\n")
		// continue
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		program := ProgramData{
			Name:  strings.TrimSpace(line), // 应用程序的名称
			APPID: "",                      // 应用程序的 APPID
		}
		programs = append(programs, program)
	}

	cmd = exec.Command("powershell", "chcp 65001; Get-StartApps | Select AppID")

	var out2 bytes.Buffer
	cmd.Stdout = &out2

	// 执行命令并捕获输出
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	// 解析 PowerShell 输出
	lines = strings.Split(out2.String(), "\n")
	for index, line := range lines[4:] { // 跳过标题行

		// 假设 Name 和 APPID 之间有空格
		// parts := strings.Fields(line)

		// print(index)
		line = strings.TrimSpace(line)

		if index < len(programs) {
			programs[index].APPID = strings.TrimSpace(line)
		}

	}

	return programs, nil
}

func main() {
	// 获取开始菜单应用
	apps, err := getStartApps()
	if err != nil {
		log.Fatal(err)
	}

	// 输出所有应用名称
	// fmt.Println("Start Menu Applications:")
	for _, app := range apps {
		fmt.Println(app)
	}
}
