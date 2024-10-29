// win.go
//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

func main() {
	// 定义要查询的注册表键
	keys := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,             // 常规程序
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, // 32位程序
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Uninstall`,   // Windows 10 应用
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Appx\AppxAllUserStore`, // Windows Store 应用
	}

	for _, keyPath := range keys {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.READ)
		if err != nil {
			log.Printf("无法打开注册表键 %s: %v\n", keyPath, err)
			continue
		}
		defer key.Close()

		// 获取子键名
		subkeys, err := key.SubKeyNames(-1)
		if err != nil {
			log.Printf("无法获取子键名: %v\n", err)
			continue
		}

		// 遍历子键
		for _, subkey := range subkeys {
			appKey, err := registry.OpenKey(key, subkey, registry.READ)
			if err != nil {
				continue // 继续下一个子键
			}
			defer appKey.Close()

			// 获取应用程序名称
			displayName, _, err := appKey.GetStringValue("DisplayName")
			if err != nil {
				continue // 没有显示名称，跳过
			}

			// 获取安装路径
			installLocation, _, err := appKey.GetStringValue("InstallLocation")
			if err != nil {
				// 尝试获取 UninstallString 作为可执行文件路径
				uninstallString, _, err := appKey.GetStringValue("UninstallString")
				if err == nil {
					fmt.Printf("Application: %s\nExecutable Path: %s\n", displayName, uninstallString)
				}
				continue // 没有安装路径，跳过
			}

			// 如果 InstallLocation 不为空，使用它作为可执行文件路径
			fmt.Printf("Application: %s\nExecutable Path: %s\n", displayName, installLocation)
		}
	}
}
