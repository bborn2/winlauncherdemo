package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

func getInstalledApps() []string {
	var apps []string
	// 定义注册表路径
	registryPaths := []struct {
		root registry.Key
		path string
	}{
		{registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Uninstall`},
	}

	for _, regPath := range registryPaths {
		k, err := registry.OpenKey(regPath.root, regPath.path, registry.READ)
		if err != nil {
			log.Printf("Failed to open registry key: %v", err)
			continue
		}
		defer k.Close()

		// 枚举子键，即每个已安装应用程序的唯一 ID
		subKeys, err := k.ReadSubKeyNames(-1)
		if err != nil {
			log.Printf("Failed to read subkeys: %v", err)
			continue
		}

		for _, subKey := range subKeys {
			appKey, err := registry.OpenKey(k, subKey, registry.READ)
			if err != nil {
				log.Printf("Failed to open app key: %v", err)
				continue
			}
			defer appKey.Close()

			// 读取应用程序的显示名称
			displayName, _, err := appKey.GetStringValue("DisplayName")
			if err == nil {
				apps = append(apps, displayName)
			}
		}
	}
	return apps
}

func main() {
	apps := getInstalledApps()
	fmt.Println("Installed applications:")
	for _, app := range apps {
		fmt.Println(app)
	}
}
