// win.go
//go:build windows
// +build windows

package main

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func GetInstalledPrograms() ([]Program, error) {
	var programs []Program
	programMap := make(map[string]string) // 用于去重

	// 定义需要遍历的注册表路径
	keys := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,             // 常规程序
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, // 32位程序
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Uninstall`,   // Windows 10 应用
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Appx\AppxAllUserStore`, // Windows Store 应用
	}

	// 遍历注册表路径
	for _, keyPath := range keys {
		programsFromKey, err := getProgramsFromRegistry(keyPath)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", keyPath, err)
			continue
		}

		// 去重处理
		for _, program := range programsFromKey {
			if _, exists := programMap[program.Name]; !exists {
				programMap[program.Name] = program.Path
				programs = append(programs, program)
			}
		}
	}

	return programs, nil
}

// 从指定注册表路径获取程序
func getProgramsFromRegistry(keyPath string) ([]Program, error) {
	var programs []Program

	// 打开 LOCAL_MACHINE 注册表键
	uninstallKey, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.READ)
	if err != nil {
		return nil, err
	}
	defer uninstallKey.Close()

	// 读取子键
	names, err := uninstallKey.ReadSubKeyNames(-1)
	if err != nil {
		return nil, err
	}

	// 遍历每个子键获取 DisplayName 和 DisplayIcon
	for _, name := range names {
		subKey, err := registry.OpenKey(uninstallKey, name, registry.READ)
		if err != nil {
			continue
		}
		defer subKey.Close()

		displayName, _, err := subKey.GetStringValue("DisplayName")
		if err != nil || displayName == "" {
			continue
		}

		displayPath, _, err := subKey.GetStringValue("DisplayIcon")
		if err != nil || displayPath == "" {
			continue
		}

		programs = append(programs, Program{Name: displayName, Path: displayPath})
	}

	return programs, nil
}
