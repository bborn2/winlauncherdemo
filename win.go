// win.go
//go:build windows
// +build windows

package main

import (
	"log"

	"golang.org/x/sys/windows/registry"
)

// GetInstalledApps 返回系统中所有已安装应用程序的名称
func GetInstalledApps() []string {
	var apps []string

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

		i := 0
		for {
			subkeyName, err := key.EnumKey(i)
			if err != nil {
				break // 没有更多子键，结束循环
			}
			i++

			appKey, err := registry.OpenKey(key, subkeyName, registry.READ)
			if err != nil {
				continue // 继续下一个子键
			}
			defer appKey.Close()

			// 获取应用程序显示名称
			displayName, _, err := appKey.GetStringValue("DisplayName")
			if err == nil && displayName != "" {
				apps = append(apps, displayName)
				continue
			}

			// 尝试使用其他键名来获取 Windows Store 应用的显示名称
			appName, _, err := appKey.GetStringValue("Name")
			if err == nil && appName != "" {
				apps = append(apps, appName)
			}
		}
	}
	return apps
}
