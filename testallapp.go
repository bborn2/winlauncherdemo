package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

func getStartMenuShortcuts() ([]string, error) {
	// 用户的开始菜单应用程序路径
	userStartMenuPath := filepath.Join("C:", "Users", "songkun2", "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs")
	// 系统的开始菜单应用程序路径
	systemStartMenuPath := filepath.Join("C:", "ProgramData", "Microsoft", "Windows", "Start Menu", "Programs")

	var allShortcuts []string

	// 获取用户开始菜单中的快捷方式
	userFiles, err := ioutil.ReadDir(userStartMenuPath)
	if err != nil {
		return nil, err
	}
	for _, file := range userFiles {
		if file.IsDir() || filepath.Ext(file.Name()) != ".lnk" {
			continue
		}
		allShortcuts = append(allShortcuts, filepath.Join(userStartMenuPath, file.Name()))
	}

	// 获取系统开始菜单中的快捷方式
	systemFiles, err := ioutil.ReadDir(systemStartMenuPath)
	if err != nil {
		return nil, err
	}
	for _, file := range systemFiles {
		if file.IsDir() || filepath.Ext(file.Name()) != ".lnk" {
			continue
		}
		allShortcuts = append(allShortcuts, filepath.Join(systemStartMenuPath, file.Name()))
	}

	return allShortcuts, nil
}

func main() {
	// 获取开始菜单快捷方式
	shortcuts, err := getStartMenuShortcuts()
	if err != nil {
		log.Fatal(err)
	}

	// 输出所有快捷方式
	fmt.Println("Start Menu Shortcuts:")
	for _, shortcut := range shortcuts {
		fmt.Println(shortcut)
	}
}
